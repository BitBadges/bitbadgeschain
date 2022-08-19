/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

/**
 * Defines an approval object for a specific address.
 */
export interface BadgesApproval {
  /** @format uint64 */
  address?: string;
  approvalAmounts?: BadgesBalanceObject[];
}

/**
 * Defines a balance object. The specified balance holds for all ids specified within the id ranges array.
 */
export interface BadgesBalanceObject {
  /** @format uint64 */
  balance?: string;
  id_ranges?: BadgesIdRange[];
}

/**
 * BitBadge defines a badge type. Think of this like the smart contract definition.
 */
export interface BadgesBitBadge {
  /** @format uint64 */
  id?: string;

  /**
   * uri object for the badge uri and subasset uris stored off chain. Stored in a special UriObject that attemtps to save space and avoid reused plaintext storage such as http:// and duplicate text for uri and subasset uris
   * data returned should corresponds to the Badge standard defined.
   */
  uri?: BadgesUriObject;

  /** @format byte */
  arbitraryBytes?: string;

  /** @format uint64 */
  manager?: string;

  /** @format uint64 */
  permissions?: string;

  /** FreezeRanges defines what addresses are frozen or unfrozen. If permissions.FrozenByDefault is false, this is used for frozen addresses. If true, this is used for unfrozen addresses. */
  freezeRanges?: BadgesIdRange[];

  /**
   * Starts at 0. Each subasset created will incrementally have an increasing ID #. Can't overflow.
   * @format uint64
   */
  nextSubassetId?: string;
  subassetSupplys?: BadgesBalanceObject[];

  /**
   * Default subasset supply. If == 0, we assume == 1.
   * @format uint64
   */
  defaultSubassetSupply?: string;

  /**
   * Defines what standard this badge should implement. Must obey the rules of that standard.
   * @format uint64
   */
  standard?: string;
}

/**
 * Id ranges define a range of IDs from start to end. Can be used for subbadgeIds, nonces, addresses anything. If end == 0, we assume end == start. Start must be >= end.
 */
export interface BadgesIdRange {
  /** @format uint64 */
  start?: string;

  /** @format uint64 */
  end?: string;
}

export type BadgesMsgFreezeAddressResponse = object;

export type BadgesMsgHandlePendingTransferResponse = object;

export interface BadgesMsgNewBadgeResponse {
  /** @format uint64 */
  id?: string;
}

export interface BadgesMsgNewSubBadgeResponse {
  /** @format uint64 */
  nextSubassetId?: string;
}

export type BadgesMsgPruneBalancesResponse = object;

export type BadgesMsgRequestTransferBadgeResponse = object;

export type BadgesMsgRequestTransferManagerResponse = object;

export type BadgesMsgRevokeBadgeResponse = object;

export type BadgesMsgSelfDestructBadgeResponse = object;

export type BadgesMsgSetApprovalResponse = object;

export type BadgesMsgTransferBadgeResponse = object;

export type BadgesMsgTransferManagerResponse = object;

export type BadgesMsgUpdateBytesResponse = object;

export type BadgesMsgUpdatePermissionsResponse = object;

export type BadgesMsgUpdateUrisResponse = object;

/**
 * Params defines the parameters for the module.
 */
export type BadgesParams = object;

/**
 * Defines a pending transfer object for two addresses. A pending transfer will be stored in both parties' balance objects.
 */
export interface BadgesPendingTransfer {
  /** Id ranges define a range of IDs from start to end. Can be used for subbadgeIds, nonces, addresses anything. If end == 0, we assume end == start. Start must be >= end. */
  subbadgeRange?: BadgesIdRange;

  /** @format uint64 */
  thisPendingNonce?: string;

  /** @format uint64 */
  otherPendingNonce?: string;

  /** @format uint64 */
  amount?: string;
  sent?: boolean;

  /** @format uint64 */
  to?: string;

  /** @format uint64 */
  from?: string;

  /** @format uint64 */
  approvedBy?: string;
  markedAsAccepted?: boolean;

  /** @format uint64 */
  expirationTime?: string;

  /** @format uint64 */
  cantCancelBeforeTime?: string;
}

export interface BadgesQueryGetBadgeResponse {
  /** BitBadge defines a badge type. Think of this like the smart contract definition. */
  badge?: BadgesBitBadge;
}

export interface BadgesQueryGetBalanceResponse {
  /** Defines a user balance object for a badge w/ the user's balances, nonce, pending transfers, and approvals. All subbadge IDs for a badge are handled within this object. */
  balanceInfo?: BadgesUserBalanceInfo;
}

/**
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface BadgesQueryParamsResponse {
  /** params holds all the parameters of this module. */
  params?: BadgesParams;
}

/**
 * A URI object defines a uri and subasset uri for a badge and its subbadges. Designed to save storage and avoid reused text and common patterns.
 */
export interface BadgesUriObject {
  /** @format uint64 */
  decodeScheme?: string;

  /** @format uint64 */
  scheme?: string;

  /** @format byte */
  uri?: string;

  /** The four fields below are used to convert the uri from above to the subasset URI. */
  idxRangeToRemove?: BadgesIdRange;

  /** @format uint64 */
  insertSubassetBytesIdx?: string;

  /** @format byte */
  bytesToInsert?: string;

  /** @format uint64 */
  insertIdIdx?: string;
}

/**
 * Defines a user balance object for a badge w/ the user's balances, nonce, pending transfers, and approvals. All subbadge IDs for a badge are handled within this object.
 */
export interface BadgesUserBalanceInfo {
  balanceAmounts?: BadgesBalanceObject[];

  /** @format uint64 */
  pendingNonce?: string;
  pending?: BadgesPendingTransfer[];
  approvals?: BadgesApproval[];
}

/**
* `Any` contains an arbitrary serialized protocol buffer message along with a
URL that describes the type of the serialized message.

Protobuf library provides support to pack/unpack Any values in the form
of utility functions or additional generated methods of the Any type.

Example 1: Pack and unpack a message in C++.

    Foo foo = ...;
    Any any;
    any.PackFrom(foo);
    ...
    if (any.UnpackTo(&foo)) {
      ...
    }

Example 2: Pack and unpack a message in Java.

    Foo foo = ...;
    Any any = Any.pack(foo);
    ...
    if (any.is(Foo.class)) {
      foo = any.unpack(Foo.class);
    }

 Example 3: Pack and unpack a message in Python.

    foo = Foo(...)
    any = Any()
    any.Pack(foo)
    ...
    if any.Is(Foo.DESCRIPTOR):
      any.Unpack(foo)
      ...

 Example 4: Pack and unpack a message in Go

     foo := &pb.Foo{...}
     any, err := anypb.New(foo)
     if err != nil {
       ...
     }
     ...
     foo := &pb.Foo{}
     if err := any.UnmarshalTo(foo); err != nil {
       ...
     }

The pack methods provided by protobuf library will by default use
'type.googleapis.com/full.type.name' as the type URL and the unpack
methods only use the fully qualified type name after the last '/'
in the type URL, for example "foo.bar.com/x/y.z" will yield type
name "y.z".


JSON
====
The JSON representation of an `Any` value uses the regular
representation of the deserialized, embedded message, with an
additional field `@type` which contains the type URL. Example:

    package google.profile;
    message Person {
      string first_name = 1;
      string last_name = 2;
    }

    {
      "@type": "type.googleapis.com/google.profile.Person",
      "firstName": <string>,
      "lastName": <string>
    }

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
`value` which holds the custom JSON in addition to the `@type`
field. Example (for message [google.protobuf.Duration][]):

    {
      "@type": "type.googleapis.com/google.protobuf.Duration",
      "value": "1.212s"
    }
*/
export interface ProtobufAny {
  /**
   * A URL/resource name that uniquely identifies the type of the serialized
   * protocol buffer message. This string must contain at least
   * one "/" character. The last segment of the URL's path must represent
   * the fully qualified name of the type (as in
   * `path/google.protobuf.Duration`). The name should be in a canonical form
   * (e.g., leading "." is not accepted).
   *
   * In practice, teams usually precompile into the binary all types that they
   * expect it to use in the context of Any. However, for URLs which use the
   * scheme `http`, `https`, or no scheme, one can optionally set up a type
   * server that maps type URLs to message definitions as follows:
   *
   * * If no scheme is provided, `https` is assumed.
   * * An HTTP GET on the URL must yield a [google.protobuf.Type][]
   *   value in binary format, or produce an error.
   * * Applications are allowed to cache lookup results based on the
   *   URL, or have them precompiled into a binary to avoid any
   *   lookup. Therefore, binary compatibility needs to be preserved
   *   on changes to types. (Use versioned type names to manage
   *   breaking changes.)
   *
   * Note: this functionality is not currently available in the official
   * protobuf release, and it is not used for type URLs beginning with
   * type.googleapis.com.
   *
   * Schemes other than `http`, `https` (or the empty scheme) might be
   * used with implementation specific semantics.
   */
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: keyof Omit<Body, "body" | "bodyUsed">;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType) => RequestParams | void;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType = null as any;
  private securityWorker: null | ApiConfig<SecurityDataType>["securityWorker"] = null;
  private abortControllers = new Map<CancelToken, AbortController>();

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType) => {
    this.securityData = data;
  };

  private addQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];

    return (
      encodeURIComponent(key) +
      "=" +
      encodeURIComponent(Array.isArray(value) ? value.join(",") : typeof value === "number" ? value : `${value}`)
    );
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) =>
        typeof query[key] === "object" && !Array.isArray(query[key])
          ? this.toQueryString(query[key] as QueryParamsType)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((data, key) => {
        data.append(key, input[key]);
        return data;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  private mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format = "json",
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams = (secure && this.securityWorker && this.securityWorker(this.securityData)) || {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];

    return fetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      signal: cancelToken ? this.createAbortSignal(cancelToken) : void 0,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = (null as unknown) as T;
      r.error = (null as unknown) as E;

      const data = await response[format]()
        .then((data) => {
          if (r.ok) {
            r.data = data;
          } else {
            r.error = data;
          }
          return r;
        })
        .catch((e) => {
          r.error = e;
          return r;
        });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title badges/badges.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryGetBadge
   * @summary Queries a list of GetBadge items.
   * @request GET:/trevormil/bitbadgeschain/badges/get_badge/{id}
   */
  queryGetBadge = (id: string, params: RequestParams = {}) =>
    this.request<BadgesQueryGetBadgeResponse, RpcStatus>({
      path: `/trevormil/bitbadgeschain/badges/get_badge/${id}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGetBalance
   * @summary Queries a list of GetBalance items.
   * @request GET:/trevormil/bitbadgeschain/badges/get_balance/{badgeId}/{address}
   */
  queryGetBalance = (badgeId: string, address: string, params: RequestParams = {}) =>
    this.request<BadgesQueryGetBalanceResponse, RpcStatus>({
      path: `/trevormil/bitbadgeschain/badges/get_balance/${badgeId}/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Parameters queries the parameters of the module.
   * @request GET:/trevormil/bitbadgeschain/badges/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<BadgesQueryParamsResponse, RpcStatus>({
      path: `/trevormil/bitbadgeschain/badges/params`,
      method: "GET",
      format: "json",
      ...params,
    });
}
