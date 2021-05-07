import { Reader, Writer } from "protobufjs/minimal";
import { Nym } from "../verimcosmos/nym";
import { PageRequest, PageResponse } from "../cosmos/base/query/v1beta1/pagination";
export declare const protobufPackage = "verimid.verimcosmos.verimcosmos";
/** this line is used by starport scaffolding # 3 */
export interface QueryGetNymRequest {
    id: number;
}
export interface QueryGetNymResponse {
    Nym: Nym | undefined;
}
export interface QueryAllNymRequest {
    pagination: PageRequest | undefined;
}
export interface QueryAllNymResponse {
    Nym: Nym[];
    pagination: PageResponse | undefined;
}
export declare const QueryGetNymRequest: {
    encode(message: QueryGetNymRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryGetNymRequest;
    fromJSON(object: any): QueryGetNymRequest;
    toJSON(message: QueryGetNymRequest): unknown;
    fromPartial(object: DeepPartial<QueryGetNymRequest>): QueryGetNymRequest;
};
export declare const QueryGetNymResponse: {
    encode(message: QueryGetNymResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryGetNymResponse;
    fromJSON(object: any): QueryGetNymResponse;
    toJSON(message: QueryGetNymResponse): unknown;
    fromPartial(object: DeepPartial<QueryGetNymResponse>): QueryGetNymResponse;
};
export declare const QueryAllNymRequest: {
    encode(message: QueryAllNymRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryAllNymRequest;
    fromJSON(object: any): QueryAllNymRequest;
    toJSON(message: QueryAllNymRequest): unknown;
    fromPartial(object: DeepPartial<QueryAllNymRequest>): QueryAllNymRequest;
};
export declare const QueryAllNymResponse: {
    encode(message: QueryAllNymResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryAllNymResponse;
    fromJSON(object: any): QueryAllNymResponse;
    toJSON(message: QueryAllNymResponse): unknown;
    fromPartial(object: DeepPartial<QueryAllNymResponse>): QueryAllNymResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /** this line is used by starport scaffolding # 2 */
    Nym(request: QueryGetNymRequest): Promise<QueryGetNymResponse>;
    NymAll(request: QueryAllNymRequest): Promise<QueryAllNymResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Nym(request: QueryGetNymRequest): Promise<QueryGetNymResponse>;
    NymAll(request: QueryAllNymRequest): Promise<QueryAllNymResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
