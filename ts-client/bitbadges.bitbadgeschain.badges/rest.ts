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
  /**
   * account number for the address
   * @format uint64
   */
  address?: string;

  /** approval balances for every subbadgeId */
  approvalAmounts?: BadgesBalanceObject[];
}

/**
 * Defines a balance object. The specified balance holds for all ids specified within the id ranges array.
 */
export interface BadgesBalanceObject {
  /** @format uint64 */
  balance?: string;
  idRanges?: BadgesIdRange[];
}

/**
 * BitBadge defines a badge type. Think of this like the smart contract definition.
 */
export interface BadgesBitBadge {
  /**
   * id defines the unique identifier of the Badge classification, similar to the contract address of ERC721
   * starts at 0 and increments by 1 each badge
   * @format uint64
   */
  id?: string;

  /**
   * uri object for the badge uri and subasset uris stored off chain. Stored in a special UriObject that attemtps to save space and avoid reused plaintext storage such as http:// and duplicate text for uri and subasset uris
   * data returned should corresponds to the Badge standard defined.
   */
  uri?: BadgesUriObject;

  /**
   * these bytes can be used to store anything on-chain about the badge. This can be updatable or not depending on the permissions set.
   * Max 256 bytes allowed
   */
  arbitraryBytes?: string;

  /**
   * manager address of the class; can have special permissions; is used as the reserve address for the assets
   * @format uint64
   */
  manager?: string;

  /**
   * Store permissions packed in a uint where the bits correspond to permissions from left to right; leading zeroes are applied and any future additions will be appended to the right. See types/permissions.go
   * @format uint64
   */
  permissions?: string;

  /** FreezeRanges defines what addresses are frozen or unfrozen. If permissions.FrozenByDefault is false, this is used for frozen addresses. If true, this is used for unfrozen addresses. */
  freezeRanges?: BadgesIdRange[];

  /**
   * Starts at 0. Each subasset created will incrementally have an increasing ID #. Can't overflow.
   * @format uint64
   */
  nextSubassetId?: string;

  /** Subasset supplys are stored if the subasset supply != default. Balance => SubbadgeIdRange map */
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
  /**
   * ID of created badge
   * @format uint64
   */
  id?: string;
}

export interface BadgesMsgNewSubBadgeResponse {
  /**
   * ID of next subbadgeId after creating all subbadges.
   * @format uint64
   */
  nextSubassetId?: string;
}

export type BadgesMsgPruneBalancesResponse = object;

export interface BadgesMsgRegisterAddressesResponse {
  /** Id ranges define a range of IDs from start to end. Can be used for subbadgeIds, nonces, addresses anything. If end == 0, we assume end == start. Start must be >= end. */
  registeredAddressNumbers?: BadgesIdRange;
}

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

  /**
   * This pending nonce is the nonce of the account for which this transfer is stored. Other is the other party's. Will be swapped for the other party's stored pending transfer.
   * @format uint64
   */
  thisPendingNonce?: string;

  /** @format uint64 */
  otherPendingNonce?: string;

  /** @format uint64 */
  amount?: string;

  /** Sent defines who initiated this pending transfer */
  sent?: boolean;

  /** @format uint64 */
  to?: string;

  /** @format uint64 */
  from?: string;

  /** @format uint64 */
  approvedBy?: string;

  /** For non forceful accepts, this will be true if the other party has accepted but doesn't want to pay the gas fees. */
  markedAsAccepted?: boolean;

  /**
   * Can't be accepted after expiration time. If == 0, we assume it never expires.
   * @format uint64
   */
  expirationTime?: string;

  /**
   * Can't cancel before must be less than expiration time. If == 0, we assume can cancel at any time.
   * @format uint64
   */
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
  /**
   * This will be == 0 represeting plaintext URLs for now, but in the future, if we want to add other decoding / encoding schemes like Base64, etc, we just define a new decoding scheme
   * @format uint64
   */
  decodeScheme?: string;

  /**
   * Helps to save space by not storing methods like https:// every time. If this == 0, we assume it's included in the uri itself. Else, we define a few predefined int -> scheme maps in uris.go that will prefix the uri bytes.
   * @format uint64
   */
  scheme?: string;

  /** The uri bytes to store. Will be converted to a string. To be manipulated according to other properties of this URI object. */
  uri?: string;

  /**
   * The four fields below are used to convert the uri from above to the subasset URI.
   *
   * Remove this range of characters from the uri for the subasset. Leave nil if no removal is needed.
   */
  idxRangeToRemove?: BadgesIdRange;

  /** @format uint64 */
  insertSubassetBytesIdx?: string;

  /** After removing the above range, insert these bytes at insertSubassetBytesIdx. insertSubassetBytesIdx is the idx after the range removal, not before. */
  bytesToInsert?: string;

  /**
   * This is the idx where we insert the id number of the subbadge. For example, ex.com/ and insertIdIdx == 6, the subasset URI for ID 0 would be ex.com/0
   * @format uint64
   */
  insertIdIdx?: string;
}

/**
 * Defines a user balance object for a badge w/ the user's balances, nonce, pending transfers, and approvals. All subbadge IDs for a badge are handled within this object.
 */
export interface BadgesUserBalanceInfo {
  /** The user's balance for each subbadge. */
  balanceAmounts?: BadgesBalanceObject[];

  /**
   * Nonce for pending transfers. Increments by 1 each time.
   * @format uint64
   */
  pendingNonce?: string;

  /** IDs will be sorted in order of this account's pending nonce. */
  pending?: BadgesPendingTransfer[];

  /** Approvals are sorted in order of address number. */
  approvals?: BadgesApproval[];
}

export interface BadgesWhitelistMintInfo {
  addresses?: string[];
  balanceAmounts?: BadgesBalanceObject[];
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
   * * If no scheme is provided, `https` is assumed.
   * * An HTTP GET on the URL must yield a [google.protobuf.Type][]
   *   value in binary format, or produce an error.
   * * Applications are allowed to cache lookup results based on the
   *   URL, or have them precompiled into a binary to avoid any
   *   lookup. Therefore, binary compatibility needs to be preserved
   *   on changes to types. (Use versioned type names to manage
   *   breaking changes.)
   * Note: this functionality is not currently available in the official
   * protobuf release, and it is not used for type URLs beginning with
   * type.googleapis.com.
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

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
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
   * @request GET:/bitbadges/bitbadgeschain/badges/get_badge/{id}
   */
  queryGetBadge = (id: string, params: RequestParams = {}) =>
    this.request<BadgesQueryGetBadgeResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/get_badge/${id}`,
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
   * @request GET:/bitbadges/bitbadgeschain/badges/get_balance/{badgeId}/{address}
   */
  queryGetBalance = (badgeId: string, address: string, params: RequestParams = {}) =>
    this.request<BadgesQueryGetBalanceResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/get_balance/${badgeId}/${address}`,
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
   * @request GET:/bitbadges/bitbadgeschain/badges/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<BadgesQueryParamsResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/params`,
      method: "GET",
      format: "json",
      ...params,
    });
}
