import { Reader, Writer } from "protobufjs/minimal";
export declare const protobufPackage = "verimid.verimcosmos.verimcosmos";
/** this line is used by starport scaffolding # proto/tx/message */
export interface MsgCreateNym {
    creator: string;
    alais: string;
    verkey: string;
    did: string;
    role: string;
}
export interface MsgCreateNymResponse {
    id: number;
}
export interface MsgUpdateNym {
    creator: string;
    id: number;
    alais: string;
    verkey: string;
    did: string;
    role: string;
}
export interface MsgUpdateNymResponse {
}
export interface MsgDeleteNym {
    creator: string;
    id: number;
}
export interface MsgDeleteNymResponse {
}
export declare const MsgCreateNym: {
    encode(message: MsgCreateNym, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgCreateNym;
    fromJSON(object: any): MsgCreateNym;
    toJSON(message: MsgCreateNym): unknown;
    fromPartial(object: DeepPartial<MsgCreateNym>): MsgCreateNym;
};
export declare const MsgCreateNymResponse: {
    encode(message: MsgCreateNymResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgCreateNymResponse;
    fromJSON(object: any): MsgCreateNymResponse;
    toJSON(message: MsgCreateNymResponse): unknown;
    fromPartial(object: DeepPartial<MsgCreateNymResponse>): MsgCreateNymResponse;
};
export declare const MsgUpdateNym: {
    encode(message: MsgUpdateNym, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgUpdateNym;
    fromJSON(object: any): MsgUpdateNym;
    toJSON(message: MsgUpdateNym): unknown;
    fromPartial(object: DeepPartial<MsgUpdateNym>): MsgUpdateNym;
};
export declare const MsgUpdateNymResponse: {
    encode(_: MsgUpdateNymResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgUpdateNymResponse;
    fromJSON(_: any): MsgUpdateNymResponse;
    toJSON(_: MsgUpdateNymResponse): unknown;
    fromPartial(_: DeepPartial<MsgUpdateNymResponse>): MsgUpdateNymResponse;
};
export declare const MsgDeleteNym: {
    encode(message: MsgDeleteNym, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgDeleteNym;
    fromJSON(object: any): MsgDeleteNym;
    toJSON(message: MsgDeleteNym): unknown;
    fromPartial(object: DeepPartial<MsgDeleteNym>): MsgDeleteNym;
};
export declare const MsgDeleteNymResponse: {
    encode(_: MsgDeleteNymResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgDeleteNymResponse;
    fromJSON(_: any): MsgDeleteNymResponse;
    toJSON(_: MsgDeleteNymResponse): unknown;
    fromPartial(_: DeepPartial<MsgDeleteNymResponse>): MsgDeleteNymResponse;
};
/** Msg defines the Msg service. */
export interface Msg {
    /** this line is used by starport scaffolding # proto/tx/rpc */
    CreateNym(request: MsgCreateNym): Promise<MsgCreateNymResponse>;
    UpdateNym(request: MsgUpdateNym): Promise<MsgUpdateNymResponse>;
    DeleteNym(request: MsgDeleteNym): Promise<MsgDeleteNymResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    CreateNym(request: MsgCreateNym): Promise<MsgCreateNymResponse>;
    UpdateNym(request: MsgUpdateNym): Promise<MsgUpdateNymResponse>;
    DeleteNym(request: MsgDeleteNym): Promise<MsgDeleteNymResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
