/**
 * Build the reader-facing OpenAPI 3.1 document for the BitBadges chain LCD.
 *
 * THIS FILE IS THE SOURCE OF TRUTH. The chain repo owns the published API
 * document; the docs site (bitbadges-docs `site/`) only consumes the output of
 * this script and hosts it with Scalar at `/chain-api-reference`. Run this
 * whenever `docs/static/openapi.yml` is regenerated, and commit the output at
 * `docs/openapi/openapi.json`.
 *
 * THE OUTPUT MUST NOT LIVE IN `docs/static/`. `docs/docs.go` does
 * `//go:embed static`, so every byte under that directory is compiled into the
 * `bitbadgeschaind` release binary. The swagger source there is already ~992 KB
 * and the node only ever serves the `.yml`; writing the converted document back
 * into `static/` would add roughly another megabyte to every release for no
 * runtime benefit. `docs/openapi/` sits outside the embedded tree — mirroring
 * how bitbadgesjs keeps `openapi-hosted/`. Do not move it back.
 *
 * Why it exists: the Cosmos swagger generator emits a Swagger 2.0 document with
 * exactly two tags — `Query` and `Msg` — across 280+ paths. That is unusable as
 * a reference; every operation lands in one of two buckets and the reader
 * cannot tell `x/tokenization` from `x/gamm`. This script converts the document
 * to OpenAPI 3.1 and retags every operation by the module its path belongs to.
 *
 * It also drops the paths that are not routes. buf writes every gRPC method
 * that carries no `google.api.http` annotation as a pseudo-path built from its
 * gRPC selector — `/tokenization.Msg/CreateCollection`,
 * `/gamm.v1beta1.Query/CalcJoinPoolNoSwapShares`. The LCD answers those with
 * `501 Not Implemented` (verified against `lcd.bitbadges.io`), so publishing
 * them gives the reader a playground button that cannot work. Transactions are
 * signed and broadcast, never POSTed to the LCD; the surface this document
 * describes is read-only.
 *
 * PORTABILITY CONTRACT: this file is self-contained on purpose. It imports
 * nothing but `node:` builtins — no npm dependency at all — so the identical
 * body lives in both repos and can be copied either way unchanged. Input and
 * output paths come from argv/env. Keep it that way.
 *
 * Usage:
 *   bun scripts/gen-openapi.ts [--in docs/static/openapi.yml] [--out docs/openapi/openapi.json]
 */
import fs from 'node:fs/promises';
import path from 'node:path';

export type Json = Record<string, any>;

/* ==========================================================================
   1. Swagger 2.0 -> OpenAPI 3.1
   ========================================================================== */

export const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);

/** Parameter keys that stay on the Parameter Object; everything else is schema. */
const PARAMETER_KEYS = new Set(['name', 'in', 'description', 'required', 'deprecated', 'allowEmptyValue']);

const COLLECTION_FORMATS: Record<string, { style: string; explode: boolean }> = {
  csv: { style: 'form', explode: false },
  multi: { style: 'form', explode: true },
  ssv: { style: 'spaceDelimited', explode: false },
  pipes: { style: 'pipeDelimited', explode: false },
  tsv: { style: 'form', explode: false },
};

export type ConvertReport = {
  /** `"<method> <path>"` for each body parameter turned into a requestBody. */
  requestBodies: string[];
  /** Number of `$ref`s rewritten out of the Swagger 2.0 namespaces. */
  rewrittenRefs: number;
  /** Definitions moved to `components.schemas`. */
  movedSchemas: number;
  /** Swagger constructs seen but not converted (should stay empty). */
  unhandled: string[];
};

const REF_MAP: [RegExp, string][] = [
  [/^#\/definitions\//, '#/components/schemas/'],
  [/^#\/responses\//, '#/components/responses/'],
  [/^#\/parameters\//, '#/components/parameters/'],
];

/** Rewrite every local `$ref` from the Swagger 2.0 layout to the 3.x one. */
export function rewriteRefs(node: unknown, counter = { n: 0 }): unknown {
  if (Array.isArray(node)) return node.map((item) => rewriteRefs(item, counter));
  if (typeof node !== 'object' || node === null) return node;

  const out: Json = {};
  for (const [key, value] of Object.entries(node as Json)) {
    if (key === '$ref' && typeof value === 'string') {
      const rule = REF_MAP.find(([pattern]) => pattern.test(value));
      if (rule) {
        counter.n += 1;
        out.$ref = value.replace(rule[0], rule[1]);
        continue;
      }
    }
    out[key] = rewriteRefs(value, counter);
  }
  return out;
}

/**
 * Convert a Swagger 2.0 schema to a JSON Schema 2020-12 one.
 *
 * The transform is small because the constructs that actually differ are few:
 * `x-nullable`/`nullable` becomes a type union, boolean `exclusiveMinimum`
 * becomes a numeric bound, and `type: file` becomes a binary string.
 */
export function convertSchema(node: unknown, unhandled: Set<string>): unknown {
  if (Array.isArray(node)) return node.map((item) => convertSchema(item, unhandled));
  if (typeof node !== 'object' || node === null) return node;

  const source = node as Json;
  const out: Json = {};

  for (const [key, value] of Object.entries(source)) {
    if (key === 'nullable' || key === 'x-nullable') continue;
    if (key === 'exclusiveMinimum' || key === 'exclusiveMaximum') continue;
    if (key === 'discriminator' && typeof value === 'string') {
      out.discriminator = { propertyName: value };
      continue;
    }
    if (key === 'required' && Array.isArray(value) && value.length === 0) continue;
    out[key] = convertSchema(value, unhandled);
  }

  if (out.type === 'file') {
    out.type = 'string';
    out.format = 'binary';
  }

  // OpenAPI 3.1 / JSON Schema 2020-12 replaced the boolean flags with numbers.
  for (const [flag, bound] of [
    ['exclusiveMinimum', 'minimum'],
    ['exclusiveMaximum', 'maximum'],
  ] as const) {
    if (source[flag] === true && typeof source[bound] === 'number') {
      out[flag] = source[bound];
      delete out[bound];
    } else if (typeof source[flag] === 'number') {
      out[flag] = source[flag];
    } else if (source[flag] === false) {
      // The default; nothing to carry over.
    } else if (source[flag] !== undefined) {
      unhandled.add(`${flag}: ${JSON.stringify(source[flag])}`);
    }
  }

  if (source.nullable === true || source['x-nullable'] === true) {
    if (typeof out.type === 'string') out.type = [out.type, 'null'];
    else if (Array.isArray(out.type) && !out.type.includes('null')) out.type = [...out.type, 'null'];
  }

  return out;
}

/** Split a Swagger 2.0 parameter list into 3.x parameters plus a requestBody. */
export function convertParameters(
  parameters: Json[],
  consumes: string[],
  unhandled: Set<string>,
): { parameters: Json[]; requestBody?: Json } {
  const out: Json[] = [];
  let requestBody: Json | undefined;
  const formData: Json = { type: 'object', properties: {}, required: [] as string[] };
  let sawFormData = false;

  for (const parameter of parameters) {
    if (parameter.in === 'body') {
      requestBody = {
        description: parameter.description,
        required: parameter.required ?? false,
        content: Object.fromEntries(
          consumes.map((type) => [type, { schema: convertSchema(parameter.schema ?? {}, unhandled) }]),
        ),
      };
      if (requestBody.description === undefined) delete requestBody.description;
      continue;
    }

    if (parameter.in === 'formData') {
      sawFormData = true;
      const { name, required, in: _in, collectionFormat: _cf, ...rest } = parameter;
      formData.properties[name] = convertSchema(rest, unhandled);
      if (required) formData.required.push(name);
      continue;
    }

    const converted: Json = {};
    const schema: Json = {};
    for (const [key, value] of Object.entries(parameter)) {
      if (PARAMETER_KEYS.has(key)) converted[key] = value;
      else if (key === 'collectionFormat') {
        const style = COLLECTION_FORMATS[value as string];
        if (style) Object.assign(converted, style);
        else unhandled.add(`collectionFormat: ${value}`);
      } else if (key === 'x-example') schema.example = value;
      else schema[key] = value;
    }
    // `path` parameters are always required in 3.x; the generator sometimes
    // omits the flag, which makes the document invalid.
    if (converted.in === 'path') converted.required = true;
    converted.schema = convertSchema(schema, unhandled);
    out.push(converted);
  }

  if (sawFormData) {
    if (formData.required.length === 0) delete formData.required;
    requestBody = {
      required: true,
      content: { 'application/x-www-form-urlencoded': { schema: formData } },
    };
  }

  return requestBody ? { parameters: out, requestBody } : { parameters: out };
}

/** Convert a Swagger 2.0 response into a 3.x one. */
function convertResponse(response: Json, produces: string[], unhandled: Set<string>): Json {
  const { schema, headers, examples, ...rest } = response;
  const out: Json = { ...rest };
  if (out.description === undefined) out.description = '';

  if (schema !== undefined) {
    out.content = Object.fromEntries(
      produces.map((type) => [
        type,
        {
          schema: convertSchema(schema, unhandled),
          ...(examples && examples[type] !== undefined ? { example: examples[type] } : {}),
        },
      ]),
    );
  }

  if (headers) {
    out.headers = Object.fromEntries(
      Object.entries(headers as Json).map(([name, header]) => {
        const { description, ...schemaPart } = header as Json;
        return [name, { ...(description ? { description } : {}), schema: convertSchema(schemaPart, unhandled) }];
      }),
    );
  }

  return out;
}

function convertSecurityDefinitions(definitions: Json, unhandled: Set<string>): Json {
  const out: Json = {};
  for (const [name, scheme] of Object.entries(definitions)) {
    const source = scheme as Json;
    if (source.type === 'basic') out[name] = { type: 'http', scheme: 'basic' };
    else if (source.type === 'apiKey') out[name] = { ...source };
    else if (source.type === 'oauth2') {
      const flowName = { implicit: 'implicit', password: 'password', application: 'clientCredentials', accessCode: 'authorizationCode' }[
        source.flow as string
      ];
      if (!flowName) {
        unhandled.add(`oauth2 flow: ${source.flow}`);
        continue;
      }
      out[name] = {
        type: 'oauth2',
        flows: {
          [flowName]: {
            ...(source.authorizationUrl ? { authorizationUrl: source.authorizationUrl } : {}),
            ...(source.tokenUrl ? { tokenUrl: source.tokenUrl } : {}),
            scopes: source.scopes ?? {},
          },
        },
      };
    } else unhandled.add(`securityDefinition type: ${source.type}`);
  }
  return out;
}

/** Build `servers` from `schemes` + `host` + `basePath`. */
export function serversFromSwagger(input: Json, fallback: string): Json[] {
  if (!input.host) return [{ url: fallback }];
  const schemes: string[] = Array.isArray(input.schemes) && input.schemes.length ? input.schemes : ['https'];
  const basePath = (input.basePath ?? '').replace(/\/$/, '');
  return schemes.map((scheme) => ({ url: `${scheme}://${input.host}${basePath}` }));
}

export type ConvertOptions = { fallbackServer?: string };

/** Swagger 2.0 -> OpenAPI 3.1. The input document is never mutated. */
export function swaggerToOpenApi31(input: Json, options: ConvertOptions = {}): { spec: Json; report: ConvertReport } {
  if (input.swagger !== '2.0') throw new Error(`expected a swagger 2.0 document, got ${JSON.stringify(input.swagger)}`);

  const unhandled = new Set<string>();
  const requestBodies: string[] = [];
  const globalConsumes: string[] = input.consumes?.length ? input.consumes : ['application/json'];
  const globalProduces: string[] = input.produces?.length ? input.produces : ['application/json'];

  const paths: Json = {};
  for (const [route, item] of Object.entries((input.paths ?? {}) as Json)) {
    const source = item as Json;
    const target: Json = {};

    for (const [key, value] of Object.entries(source)) {
      if (!HTTP_METHODS.has(key)) {
        if (key === 'parameters') {
          target.parameters = convertParameters(value as Json[], globalConsumes, unhandled).parameters;
        } else target[key] = value;
        continue;
      }

      const operation = value as Json;
      const consumes = operation.consumes?.length ? operation.consumes : globalConsumes;
      const produces = operation.produces?.length ? operation.produces : globalProduces;
      const { consumes: _c, produces: _p, parameters, responses, schemes, ...rest } = operation;

      const converted: Json = { ...rest };
      if (parameters) {
        const split = convertParameters(parameters as Json[], consumes, unhandled);
        if (split.parameters.length) converted.parameters = split.parameters;
        if (split.requestBody) {
          converted.requestBody = split.requestBody;
          requestBodies.push(`${key} ${route}`);
        }
      }
      converted.responses = Object.fromEntries(
        Object.entries((responses ?? {}) as Json).map(([code, response]) => [
          code,
          convertResponse(response as Json, produces, unhandled),
        ]),
      );
      target[key] = converted;
    }

    paths[route] = target;
  }

  const schemas = Object.fromEntries(
    Object.entries((input.definitions ?? {}) as Json).map(([name, schema]) => [name, convertSchema(schema, unhandled)]),
  );

  const components: Json = { schemas };
  if (input.responses) {
    components.responses = Object.fromEntries(
      Object.entries(input.responses as Json).map(([name, response]) => [
        name,
        convertResponse(response as Json, globalProduces, unhandled),
      ]),
    );
  }
  if (input.parameters) {
    components.parameters = Object.fromEntries(
      Object.entries(input.parameters as Json).map(([name, parameter]) => [
        name,
        convertParameters([parameter as Json], globalConsumes, unhandled).parameters[0],
      ]),
    );
  }
  if (input.securityDefinitions) {
    components.securitySchemes = convertSecurityDefinitions(input.securityDefinitions as Json, unhandled);
  }

  const {
    swagger: _s,
    info,
    host: _h,
    basePath: _b,
    schemes: _sch,
    consumes: _c,
    produces: _p,
    paths: _paths,
    definitions: _d,
    responses: _r,
    parameters: _param,
    securityDefinitions: _sd,
    id: _id,
    ...passthrough
  } = input;

  const counter = { n: 0 };
  const spec = rewriteRefs(
    {
      openapi: '3.1.0',
      info: { ...(info ?? {}) },
      servers: serversFromSwagger(input, options.fallbackServer ?? '/'),
      ...passthrough,
      paths,
      components,
    },
    counter,
  ) as Json;

  return {
    spec,
    report: {
      requestBodies: requestBodies.sort(),
      rewrittenRefs: counter.n,
      movedSchemas: Object.keys(schemas).length,
      unhandled: [...unhandled].sort(),
    },
  };
}

/* ==========================================================================
   2. Prune everything that is not an LCD route
   ========================================================================== */

/**
 * A gRPC selector standing in for an HTTP route: one leading segment of the
 * form `<package>.<Service>` followed by `/<Method>`, e.g.
 * `/tokenization.Msg/CreateCollection` or
 * `/gamm.v1beta1.Query/CalcJoinPoolNoSwapShares`. Real gateway routes are
 * multi-segment and lowercase (`/osmosis/gamm/v1beta1/pools`), so nothing
 * genuine matches this shape.
 */
export const PSEUDO_PATH = /^\/[a-z][a-z0-9_]*(?:\.[a-z0-9_]+)*\.[A-Z][A-Za-z0-9]*\/[A-Za-z0-9_]+$/;

export function isPseudoPath(route: string): boolean {
  return PSEUDO_PATH.test(route);
}

export type PruneReport = {
  /** Paths removed because they are gRPC selectors, not routes. */
  pseudoPaths: string[];
  /** `"<method> <path>"` for each non-GET operation removed. */
  writeOperations: string[];
  /** Schemas left unreachable once the pseudo-paths were gone. */
  orphanSchemas: number;
};

/**
 * Remove the paths the LCD does not serve.
 *
 * Two kinds: gRPC selector pseudo-paths (see `PSEUDO_PATH`), and any non-GET
 * operation. The gateway implements queries only — a `POST` to it returns
 * `501`, whether it is a selector path or an annotated one like
 * `/cosmos/evm/vm/v1/ethereum_tx`. Messages are documented as payload shapes at
 * `/token-standard/messages`, not as callable endpoints here.
 */
export function pruneNonRoutes(input: Json): { spec: Json; report: PruneReport } {
  const spec = structuredClone(input);
  const pseudoPaths: string[] = [];
  const writeOperations: string[] = [];

  for (const [route, item] of Object.entries((spec.paths ?? {}) as Json)) {
    if (isPseudoPath(route)) {
      delete spec.paths[route];
      pseudoPaths.push(route);
      continue;
    }
    for (const method of Object.keys(item as Json)) {
      if (!HTTP_METHODS.has(method) || method === 'get') continue;
      delete (item as Json)[method];
      writeOperations.push(`${method} ${route}`);
    }
    if (!Object.keys(item as Json).some((key) => HTTP_METHODS.has(key))) delete spec.paths[route];
  }

  // The dropped paths carried most of the document's schemas with them (every
  // `Msg*` request and response). Leaving them behind would keep a ~1 MB
  // payload and fill Scalar's model list with types no operation uses.
  const schemas: Json = spec.components?.schemas ?? {};
  const reachable = new Set<string>();
  const queue = [...collectRefs(spec.paths ?? {}, new Set())];
  while (queue.length) {
    const name = queue.pop()!;
    if (reachable.has(name) || !(name in schemas)) continue;
    reachable.add(name);
    queue.push(...collectRefs(schemas[name], new Set()));
  }
  let orphanSchemas = 0;
  for (const name of Object.keys(schemas)) {
    if (reachable.has(name)) continue;
    delete schemas[name];
    orphanSchemas += 1;
  }

  return {
    spec,
    report: { pseudoPaths: pseudoPaths.sort(), writeOperations: writeOperations.sort(), orphanSchemas },
  };
}

/* ==========================================================================
   3. Retag by module
   ========================================================================== */

export type ModuleTag = { name: string; description: string };

/** Path prefix -> tag. First match wins. */
const MODULE_RULES: [RegExp, string][] = [
  [/^\/bitbadges\/bitbadgeschain\/tokenization\//, 'Tokenization'],
  [/^\/tokenization\./, 'Tokenization'],
  [/^\/bitbadges\/bitbadgeschain\/managersplitter/, 'Manager splitter'],
  [/^\/managersplitter\./, 'Manager splitter'],
  [/^\/bitbadges\/bitbadgeschain\/sendmanager/, 'Send manager'],
  [/^\/sendmanager\./, 'Send manager'],
  [/^\/osmosis\/poolmanager\//, 'Pool manager'],
  [/^\/poolmanager\./, 'Pool manager'],
  [/^\/osmosis\/gamm\//, 'GAMM'],
  [/^\/gamm\./, 'GAMM'],
  [/^\/bitbadges\/bitbadgeschain\/ibcratelimit/, 'IBC rate limit'],
  [/^\/ibcratelimit\./, 'IBC rate limit'],
  [/^\/cosmos\/evm\//, 'EVM'],
  [/^\/cosmos\.evm\./, 'EVM'],
  [/^\/ethermint[./]/, 'EVM'],
  [/^\/ibc[./]/, 'IBC'],
  [/^\/cosmos[./]/, 'Cosmos SDK'],
  [/^\/capability[./]/, 'Cosmos SDK'],
];

/**
 * Sidebar order. BitBadges modules first, then the standard Cosmos/IBC/EVM
 * surface, then anything unattributable.
 */
export const TAG_ORDER: ModuleTag[] = [
  {
    name: 'Tokenization',
    description:
      '`x/tokenization` — the BitBadges token standard: collections, balances, address lists, approval trackers and dynamic stores, read straight from node state.',
  },
  {
    name: 'GAMM',
    description: '`x/gamm` — the AMM pools: pool state, spot prices, share math and swap estimates.',
  },
  {
    name: 'Pool manager',
    description: '`x/poolmanager` — routing across pools, taker fee agreements and multi-hop swap estimation.',
  },
  {
    name: 'Send manager',
    description: '`x/sendmanager` — alias routing for sends, and the balances behind an alias denom.',
  },
  {
    name: 'Manager splitter',
    description: '`x/managersplitter` — splitting a collection manager across several addresses.',
  },
  {
    name: 'IBC rate limit',
    description: '`x/ibc-rate-limit` — the outbound and inbound IBC transfer rate limits and their parameters.',
  },
  { name: 'EVM', description: 'The Cosmos EVM module — VM parameters and state. Contract calls go to the EVM JSON-RPC, not here.' },
  { name: 'Cosmos SDK', description: 'Standard Cosmos SDK module queries (auth, bank, staking, gov and friends), unchanged from upstream.' },
  { name: 'IBC', description: 'Standard IBC queries (clients, connections, channels, transfer), unchanged from upstream.' },
  { name: 'Other', description: 'Queries this document could not attribute to a module.' },
];

export type RetagReport = {
  /** Tag name -> operation count, in sidebar order. */
  counts: [string, number][];
};

/** The module a path belongs to. */
export function moduleForPath(route: string): string {
  for (const [pattern, name] of MODULE_RULES) if (pattern.test(route)) return name;
  return 'Other';
}

/** Replace the source document's `Query`/`Msg` tags with one module tag per operation. */
export function retagByModule(input: Json): { spec: Json; report: RetagReport } {
  const spec = structuredClone(input);
  const counts = new Map<string, number>();

  for (const [route, item] of Object.entries((spec.paths ?? {}) as Json)) {
    const tag = moduleForPath(route);
    for (const [method, operation] of Object.entries(item as Json)) {
      if (!HTTP_METHODS.has(method)) continue;
      const op = operation as Json;
      op.tags = [tag];
      counts.set(tag, (counts.get(tag) ?? 0) + 1);
    }
  }

  const used = TAG_ORDER.filter((tag) => counts.has(tag.name));
  spec.tags = used;

  // Sidebar order follows `tags`, but Scalar also walks `paths` in insertion
  // order within a tag, so reorder paths to match the tag order too.
  const order = new Map(used.map((tag, index) => [tag.name, index]));
  const routes = Object.keys(spec.paths ?? {});
  routes.sort((a, b) => {
    const delta = (order.get(moduleForPath(a)) ?? 99) - (order.get(moduleForPath(b)) ?? 99);
    return delta !== 0 ? delta : a.localeCompare(b);
  });
  spec.paths = Object.fromEntries(routes.map((route) => [route, spec.paths[route]]));

  return { spec, report: { counts: used.map((tag) => [tag.name, counts.get(tag.name) ?? 0]) } };
}

/* ==========================================================================
   4. Short summaries

   Scalar labels each sidebar entry with the operation's `summary`, falling back
   to the raw path. The proto comments the generator lifts into `summary` are
   sentences and sometimes paragraphs — "SpotPrice defines a gRPC query handler
   that returns the spot price given a base denomination and a quote
   denomination." — which turns the sidebar into a wall of prose. Every
   operation gets a short imperative label here; the proto comment moves to
   `description`, where a reference renderer expects it.
   ========================================================================== */

/** Longest acceptable sidebar label. */
export const SUMMARY_MAX = 40;

/** Method names whose mechanical humanization reads badly or runs long. */
const SUMMARY_OVERRIDES: Record<string, string> = {
  AllRegisteredAlloyedPools: 'List alloyed pools',
  AllTakerFeeShareAccumulators: 'List taker fee accumulators',
  AllTakerFeeShareAgreements: 'List taker fee agreements',
  CalcExitPoolCoinsFromShares: 'Calculate exit pool coins',
  CalcJoinPoolShares: 'Calculate join pool shares',
  EstimateSinglePoolSwapExactAmountIn: 'Estimate a single-pool swap in',
  EstimateSinglePoolSwapExactAmountOut: 'Estimate a single-pool swap out',
  EstimateSwapExactAmountIn: 'Estimate a swap in',
  EstimateSwapExactAmountInWithPrimitiveTypes: 'Estimate a swap in (primitives)',
  EstimateSwapExactAmountOut: 'Estimate a swap out',
  EstimateSwapExactAmountOutWithPrimitiveTypes: 'Estimate a swap out (primitives)',
  EstimateTradeBasedOnPriceImpact: 'Estimate a trade',
  GetAllReservedProtocolAddresses: 'List reserved addresses',
  GetBalanceForToken: 'Get a token balance',
  IsAddressReservedProtocol: 'Check a reserved address',
  Params: 'Get module params',
  Parameters: 'Get module params',
  PoolsWithFilter: 'Filter pools',
  RegisteredAlloyedPoolFromDenom: 'Get an alloyed pool by denom',
  RegisteredAlloyedPoolFromPoolId: 'Get an alloyed pool by ID',
  TakerFeeShareAgreementFromDenom: 'Get a taker fee agreement',
  TakerFeeShareDenomsToAccruedValue: 'Get accrued taker fee value',
  TotalVolumeForPool: 'Get total pool volume',
};

/** Words that must keep their casing when a method name is humanized. */
const ACRONYMS = new Set(['ETH', 'EVM', 'IBC', 'ID', 'AMM', 'TWAP', 'ERC', 'LP', 'RPC', 'URI']);

/** Leading word -> verb, and whether the word itself is consumed. */
const VERB_RULES: [RegExp, string, boolean][] = [
  [/^Get$/, 'Get', true],
  [/^All$/, 'List', true],
  [/^List$/, 'List', true],
  [/^Num$/, 'Count', true],
  [/^Calc$/, 'Calculate', true],
  [/^Estimate$/, 'Estimate', true],
];

/** Split a Go/proto method name into words, keeping acronym runs together. */
export function splitWords(name: string): string[] {
  return name.match(/[A-Z]+(?![a-z])|[A-Z][a-z0-9]*|[a-z0-9]+/g) ?? [];
}

const isPlural = (word: string) => /s$/i.test(word) && !/(ss|us|is)$/i.test(word);

/** Leading nouns that read as mass quantities — "total liquidity", not "a total liquidity". */
const MASS_NOUNS = new Set(['total', 'module', 'max', 'min', 'average']);

/** The bare method name behind a gateway operationId. */
export function methodName(operationId: string | undefined, route: string): string {
  const raw = (operationId ?? '').split('_').pop() ?? '';
  const stripped = raw.replace(/Mixin\d+$/, '');
  if (stripped) return stripped;
  const last = route.split('/').filter(Boolean).pop() ?? route;
  return last.replace(/[{}]/g, '');
}

/**
 * A short imperative label for one operation — "Get a collection", "List
 * pools". Derived from the method name, disambiguated by the route's API
 * version when the chain publishes two of the same query.
 */
export function shortSummary(operationId: string | undefined, route: string): string {
  const name = methodName(operationId, route);
  const suffix = /\/v(\d+)\//.test(route) && !/^v\d+$/.test(name) ? ` (v${route.match(/\/v(\d+)\//)![1]})` : '';

  const override = SUMMARY_OVERRIDES[name.replace(/V\d+$/, '')];
  if (override) return `${override}${suffix}`;

  const words = splitWords(name.replace(/V\d+$/, ''));
  if (!words.length) return `Get ${route}`.slice(0, SUMMARY_MAX);

  let verb = 'Get';
  const rule = VERB_RULES.find(([pattern]) => pattern.test(words[0]));
  if (rule) {
    verb = rule[1];
    if (rule[2]) words.shift();
  }
  // A bare plural noun ("Pools") is a listing; a compound one ("TotalShares")
  // is a single value that happens to be plural.
  if (!rule && words.length === 1 && isPlural(words[0])) verb = 'List';
  if (!words.length) return `${verb} module params${suffix}`;

  const noun = words.map((word) => (ACRONYMS.has(word.toUpperCase()) && word.length <= 4 && word === word.toUpperCase() ? word : word.toLowerCase())).join(' ');
  const takesArticle =
    verb === 'Get' && !isPlural(words[words.length - 1]) && !MASS_NOUNS.has(words[0].toLowerCase());
  const article = takesArticle ? (/^[aeiou]/i.test(noun) ? 'an ' : 'a ') : '';
  return `${verb} ${article}${noun}${suffix}`;
}

export type SummaryReport = {
  /** `"<label>  <-  <path>"` for every label longer than `SUMMARY_MAX`. */
  overlong: string[];
  /** Operations whose original prose summary was moved to `description`. */
  movedToDescription: number;
};

/** Give every operation a short label, demoting its proto prose to `description`. */
export function summarizeOperations(input: Json): { spec: Json; report: SummaryReport } {
  const spec = structuredClone(input);
  const overlong: string[] = [];
  let movedToDescription = 0;

  for (const [route, item] of Object.entries((spec.paths ?? {}) as Json)) {
    for (const [method, operation] of Object.entries(item as Json)) {
      if (!HTTP_METHODS.has(method)) continue;
      const op = operation as Json;
      const prose = typeof op.summary === 'string' ? op.summary.replace(/\s+/g, ' ').trim() : '';

      // The proto comment is the real explanation; keep it, just not in the
      // sidebar. Single-word leftovers ("pools", "params") explain nothing.
      if (!op.description && prose.split(' ').length >= 3) {
        op.description = prose;
        movedToDescription += 1;
      }

      op.summary = shortSummary(op.operationId, route);
      if (op.summary.length > SUMMARY_MAX) overlong.push(`${op.summary}  <-  ${route}`);
    }
  }

  return { spec, report: { overlong: overlong.sort(), movedToDescription } };
}

/* ==========================================================================
   5. Sanitize — the defects that blank a strict renderer

   Scalar dereferences the whole document up front. A `$ref` that points at
   nothing, or a schema cycle, takes the page down with it rather than
   degrading. Both are stubbed here so one defect in the generator's output
   cannot cost the reader the entire reference.
   ========================================================================== */

const SCHEMA_PREFIX = '#/components/schemas/';

export type SanitizeReport = {
  /** `"<owner> -> <target>"` for each schema cycle broken. */
  cutCycles: string[];
  /** Schema names referenced but never defined. */
  stubbedRefs: string[];
  /** Path items left with no operation, removed. */
  droppedPaths: string[];
};

/** Every schema name reachable through a local `$ref`. */
function collectRefs(node: unknown, found: Set<string>): Set<string> {
  if (Array.isArray(node)) {
    for (const item of node) collectRefs(item, found);
    return found;
  }
  if (typeof node !== 'object' || node === null) return found;
  for (const [key, value] of Object.entries(node as Json)) {
    if (key === '$ref' && typeof value === 'string' && value.startsWith(SCHEMA_PREFIX)) {
      found.add(value.slice(SCHEMA_PREFIX.length));
    } else collectRefs(value, found);
  }
  return found;
}

/** Replace every `$ref` to one of `targets` with a titled stub. */
function stubRefs(node: unknown, owner: string, targets: Set<string>, cut: string[]): unknown {
  if (Array.isArray(node)) return node.map((item) => stubRefs(item, owner, targets, cut));
  if (typeof node !== 'object' || node === null) return node;

  const source = node as Json;
  if (typeof source.$ref === 'string' && source.$ref.startsWith(SCHEMA_PREFIX)) {
    const name = source.$ref.slice(SCHEMA_PREFIX.length);
    if (targets.has(name)) {
      cut.push(`${owner} -> ${name}`);
      return { title: name, description: `Recursive reference to \`${name}\`.` };
    }
  }

  const out: Json = {};
  for (const [key, value] of Object.entries(source)) out[key] = stubRefs(value, owner, targets, cut);
  return out;
}

/**
 * Break every cycle in the schema reference graph.
 *
 * A depth-first walk marks back-edges (a `$ref` reaching a schema still on the
 * stack); those edges — direct self-reference and mutual recursion alike — are
 * replaced with a titled stub. Anything else is left untouched.
 */
export function sanitizeChainSpec(input: Json): { spec: Json; report: SanitizeReport } {
  const spec = structuredClone(input);
  const schemas: Json = spec.components?.schemas ?? {};

  const edges = new Map<string, string[]>();
  for (const [name, schema] of Object.entries(schemas)) edges.set(name, [...collectRefs(schema, new Set())]);

  const state = new Map<string, 'open' | 'done'>();
  const backEdges = new Map<string, Set<string>>();
  const visit = (name: string) => {
    state.set(name, 'open');
    for (const target of edges.get(name) ?? []) {
      if (!(target in schemas)) continue;
      const seen = state.get(target);
      if (seen === 'open') {
        if (!backEdges.has(name)) backEdges.set(name, new Set());
        backEdges.get(name)!.add(target);
      } else if (seen === undefined) visit(target);
    }
    state.set(name, 'done');
  };
  for (const name of edges.keys()) if (!state.has(name)) visit(name);

  const cutCycles: string[] = [];
  for (const [owner, targets] of backEdges) schemas[owner] = stubRefs(schemas[owner], owner, targets, cutCycles);

  const stubbedRefs: string[] = [];
  for (const name of collectRefs(spec, new Set())) {
    if (name in schemas) continue;
    schemas[name] = {
      title: name,
      description: `\`${name}\` is referenced by this API but is not defined in the source document.`,
    };
    stubbedRefs.push(name);
  }
  if (Object.keys(schemas).length) {
    spec.components ??= {};
    spec.components.schemas = schemas;
  }

  const droppedPaths: string[] = [];
  for (const [route, item] of Object.entries((spec.paths ?? {}) as Json)) {
    if (!Object.keys(item as Json).some((key) => HTTP_METHODS.has(key))) {
      delete spec.paths[route];
      droppedPaths.push(route);
    }
  }

  return { spec, report: { cutCycles: cutCycles.sort(), stubbedRefs: stubbedRefs.sort(), droppedPaths: droppedPaths.sort() } };
}

/* ==========================================================================
   6. The document the reader sees
   ========================================================================== */

export const CHAIN_SERVER = 'https://lcd.bitbadges.io';

export const CHAIN_TITLE = 'BitBadges Chain API';

const CHAIN_DESCRIPTION = `The **chain LCD** — the REST surface a BitBadges node serves through the Cosmos gRPC-gateway, live at \`${CHAIN_SERVER}\`. It reads consensus state directly from a node: no indexing, no API key, no account.

## The LCD is read-only

Every route here is a \`GET\` you can run from this page against mainnet. There are no write endpoints: the gateway serves queries and answers anything else with \`501 Not Implemented\`.

Transactions never touch the LCD. A message is signed and then broadcast through the chain's RPC endpoint, or — far more simply — through the SDK. The message payloads themselves are documented at [Token Standard → Messages](/token-standard/messages); signing and broadcasting are covered in [SDK → Transactions](/sdk/transactions).

## How this differs from the BitBadges API

| | Chain API (this page) | [BitBadges API](/api-reference) |
| --- | --- | --- |
| Serves | \`https://lcd.bitbadges.io\` | \`https://api.bitbadges.io\` |
| Source of truth | A node's own state, current block | The indexer, built from chain history |
| Auth | None | API key |
| Good for | Exact on-chain values, params, pool math | Metadata, activity, claims, search, sign-in |

Use the chain API when you need the value the chain itself would return. Use the BitBadges API when you need anything the chain does not store — off-chain metadata, activity history, claims, or a query the chain has no index for.

## Documented queries

The routes here are generated from the chain's proto definitions, so each one carries only the explanation its proto comment carries. The hand-written explanation of each query — arguments, response shape, worked examples — lives at [Token Standard → Queries](/token-standard/queries), and the modules around the standard at [Chain → Modules](/chain/modules).`;

export type BuildOptions = {
  server?: string;
  title?: string;
  description?: string;
};

export type BuildReport = {
  convert: ConvertReport;
  prune: PruneReport;
  retag: RetagReport;
  summary: SummaryReport;
  sanitize: SanitizeReport;
  paths: number;
  operations: number;
  schemas: number;
};

/** Convert, prune, retag, summarize, sanitize, and stamp the reader-facing metadata. */
export function buildChainOpenApi(source: Json, options: BuildOptions = {}): { spec: Json; report: BuildReport } {
  const server = options.server ?? CHAIN_SERVER;
  const converted = swaggerToOpenApi31(source, { fallbackServer: server });
  const pruned = pruneNonRoutes(converted.spec);
  const retagged = retagByModule(pruned.spec);
  const summarized = summarizeOperations(retagged.spec);
  const sanitized = sanitizeChainSpec(summarized.spec);
  const spec = sanitized.spec;

  const version = source.info?.version;
  spec.info = {
    title: options.title ?? CHAIN_TITLE,
    version: !version || version === 'version not set' ? 'latest' : version,
    description: options.description ?? CHAIN_DESCRIPTION,
  };
  spec.servers = [{ url: server, description: 'BitBadges mainnet LCD' }];

  const operations = Object.values(spec.paths ?? {}).reduce(
    (total: number, item: any) => total + Object.keys(item).filter((key) => HTTP_METHODS.has(key)).length,
    0,
  );

  return {
    spec,
    report: {
      convert: converted.report,
      prune: pruned.report,
      retag: retagged.report,
      summary: summarized.report,
      sanitize: sanitized.report,
      paths: Object.keys(spec.paths ?? {}).length,
      operations,
      schemas: Object.keys(spec.components?.schemas ?? {}).length,
    },
  };
}

/* ==========================================================================
   7. CLI
   ========================================================================== */

/** Parse `--in <path>` and repeated `--out <path>` flags. */
export function parseArgs(argv: string[]): { input?: string; outputs: string[] } {
  const outputs: string[] = [];
  let input: string | undefined;
  for (let i = 0; i < argv.length; i += 1) {
    if (argv[i] === '--in') input = argv[++i];
    else if (argv[i] === '--out') outputs.push(argv[++i]);
  }
  return { input, outputs };
}

/**
 * Read a chain swagger/openapi document.
 *
 * The chain publishes `openapi.yml`, but the file is JSON — the extension is a
 * lie the Cosmos swagger tooling has told for years. Parsing it as JSON keeps
 * this script dependency-free; a real YAML file is rejected loudly rather than
 * half-parsed.
 */
export async function readSpec(file: string): Promise<Json> {
  const raw = await fs.readFile(file, 'utf8');
  try {
    return JSON.parse(raw);
  } catch (error) {
    throw new Error(
      `${file} is not JSON. The chain's docs/static/openapi.yml is JSON despite the extension; if it has become real YAML, convert it before running this script. (${(error as Error).message})`,
    );
  }
}

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));
  const chainDir = path.resolve(process.env.BITBADGESCHAIN_DIR?.trim() || '.');

  const usePrebuilt = false;
  const input = args.input ?? path.join(chainDir, 'docs/static/openapi.yml');

  const outputs = args.outputs.length
    ? args.outputs
    : (process.env.CHAIN_OPENAPI_OUT?.split(',').map((s) => s.trim()).filter(Boolean) ?? []);
  // Outside `docs/static/`: that tree is `//go:embed`ed into the node binary.
  if (outputs.length === 0) outputs.push(path.join(chainDir, 'docs/openapi/openapi.json'));

  const source = await readSpec(input);
  const { spec, report } =
    source.swagger === '2.0'
      ? buildChainOpenApi(source)
      : { spec: source, report: null as BuildReport | null };

  for (const out of outputs) {
    await fs.mkdir(path.dirname(out), { recursive: true });
    await fs.writeFile(out, `${JSON.stringify(spec, null, 0)}\n`);
  }

  const where = outputs.map((out) => path.relative(process.cwd(), out)).join(', ');
  if (!report) {
    console.log(`chain-openapi: copied prebuilt OpenAPI ${spec.openapi} from ${input} -> ${where}`);
    return;
  }

  const tags = report.retag.counts.map(([name, count]) => `${name} ${count}`).join(', ');
  console.log(
    `chain-openapi: ${input}${usePrebuilt ? ' (prebuilt)' : ' (swagger 2.0)'} -> ${where}\n` +
      `chain-openapi: ${report.paths} paths, ${report.operations} operations, ${report.schemas} schemas, ${report.retag.counts.length} tags\n` +
      `chain-openapi: tags — ${tags}\n` +
      `chain-openapi: dropped ${report.prune.pseudoPaths.length} gRPC pseudo-path(s) and ${report.prune.writeOperations.length} non-GET operation(s) — the LCD answers both with 501; plus ${report.prune.orphanSchemas} schema(s) left unreachable\n` +
      `chain-openapi: converted ${report.convert.requestBodies.length} body parameter(s) to requestBody, rewrote ${report.convert.rewrittenRefs} $ref(s), moved ${report.convert.movedSchemas} definition(s) to components.schemas, demoted ${report.summary.movedToDescription} prose summary(ies) to description\n` +
      `chain-openapi: fixed — ${report.sanitize.cutCycles.length} schema cycle(s), ${report.sanitize.stubbedRefs.length} dangling $ref(s), ${report.sanitize.droppedPaths.length} empty path item(s)` +
      (report.summary.overlong.length ? `\nchain-openapi: OVERLONG SUMMARIES — ${report.summary.overlong.join('; ')}` : '') +
      (report.convert.unhandled.length ? `\nchain-openapi: UNHANDLED — ${report.convert.unhandled.join('; ')}` : ''),
  );
}

if (import.meta.main) await main();
