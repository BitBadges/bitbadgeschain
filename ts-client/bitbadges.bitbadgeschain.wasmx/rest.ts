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

export interface ProtobufAny {
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
 * GenesisState defines the wasmx module's genesis state.
 */
export interface WasmxGenesisState {
  /** params defines all the parameters of related to wasmx. */
  params?: WasmxParams;
}

/**
 * MsgExecuteContractCompatResponse returns execution result data.
 */
export interface WasmxMsgExecuteContractCompatResponse {
  /**
   * Data contains bytes to returned from the contract
   * @format byte
   */
  data?: string;
}

export interface WasmxParams {
  /** Set the status to active to indicate that contracts can be executed in begin blocker. */
  is_execution_enabled?: boolean;

  /**
   * Maximum aggregate total gas to be used for the contract executions in the BeginBlocker.
   * @format uint64
   */
  max_begin_block_total_gas?: string;

  /**
   * the maximum gas limit each individual contract can consume in the BeginBlocker.
   * @format uint64
   */
  max_contract_gas_limit?: string;

  /**
   * min_gas_price defines the minimum gas price the contracts must pay to be executed in the BeginBlocker.
   * @format uint64
   */
  min_gas_price?: string;
}

/**
 * QueryModuleStateResponse is the response type for the Query/WasmxModuleState RPC method.
 */
export interface WasmxQueryModuleStateResponse {
  /** GenesisState defines the wasmx module's genesis state. */
  state?: WasmxGenesisState;
}

/**
 * QueryWasmxParamsRequest is the response type for the Query/WasmxParams RPC method.
 */
export interface WasmxQueryWasmxParamsResponse {
  params?: WasmxParams;
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
 * @title wasmx/genesis.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryWasmxModuleState
   * @summary Retrieves the entire wasmx module's state
   * @request GET:/bitbadges/bitbadgeschain/wasmx/v1/module_state
   */
  queryWasmxModuleState = (params: RequestParams = {}) =>
    this.request<WasmxQueryModuleStateResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/wasmx/v1/module_state`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryWasmxParams
   * @summary Retrieves wasmx params
   * @request GET:/bitbadges/bitbadgeschain/wasmx/v1/params
   */
  queryWasmxParams = (params: RequestParams = {}) =>
    this.request<WasmxQueryWasmxParamsResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/wasmx/v1/params`,
      method: "GET",
      format: "json",
      ...params,
    });
}
