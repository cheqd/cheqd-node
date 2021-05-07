import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "verimid.verimcosmos.verimcosmos";
export interface Nym {
    creator: string;
    id: number;
    alias: string;
    verkey: string;
    did: string;
    role: string;
}
export declare const Nym: {
    encode(message: Nym, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Nym;
    fromJSON(object: any): Nym;
    toJSON(message: Nym): unknown;
    fromPartial(object: DeepPartial<Nym>): Nym;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
