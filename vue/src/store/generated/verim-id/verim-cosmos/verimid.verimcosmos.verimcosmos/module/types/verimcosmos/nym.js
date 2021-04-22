/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
export const protobufPackage = "verimid.verimcosmos.verimcosmos";
const baseNym = {
    creator: "",
    id: 0,
    alais: "",
    verkey: "",
    did: "",
    role: "",
};
export const Nym = {
    encode(message, writer = Writer.create()) {
        if (message.creator !== "") {
            writer.uint32(10).string(message.creator);
        }
        if (message.id !== 0) {
            writer.uint32(16).uint64(message.id);
        }
        if (message.alais !== "") {
            writer.uint32(26).string(message.alais);
        }
        if (message.verkey !== "") {
            writer.uint32(34).string(message.verkey);
        }
        if (message.did !== "") {
            writer.uint32(42).string(message.did);
        }
        if (message.role !== "") {
            writer.uint32(50).string(message.role);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseNym };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.creator = reader.string();
                    break;
                case 2:
                    message.id = longToNumber(reader.uint64());
                    break;
                case 3:
                    message.alais = reader.string();
                    break;
                case 4:
                    message.verkey = reader.string();
                    break;
                case 5:
                    message.did = reader.string();
                    break;
                case 6:
                    message.role = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseNym };
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = String(object.creator);
        }
        else {
            message.creator = "";
        }
        if (object.id !== undefined && object.id !== null) {
            message.id = Number(object.id);
        }
        else {
            message.id = 0;
        }
        if (object.alais !== undefined && object.alais !== null) {
            message.alais = String(object.alais);
        }
        else {
            message.alais = "";
        }
        if (object.verkey !== undefined && object.verkey !== null) {
            message.verkey = String(object.verkey);
        }
        else {
            message.verkey = "";
        }
        if (object.did !== undefined && object.did !== null) {
            message.did = String(object.did);
        }
        else {
            message.did = "";
        }
        if (object.role !== undefined && object.role !== null) {
            message.role = String(object.role);
        }
        else {
            message.role = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.creator !== undefined && (obj.creator = message.creator);
        message.id !== undefined && (obj.id = message.id);
        message.alais !== undefined && (obj.alais = message.alais);
        message.verkey !== undefined && (obj.verkey = message.verkey);
        message.did !== undefined && (obj.did = message.did);
        message.role !== undefined && (obj.role = message.role);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseNym };
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = object.creator;
        }
        else {
            message.creator = "";
        }
        if (object.id !== undefined && object.id !== null) {
            message.id = object.id;
        }
        else {
            message.id = 0;
        }
        if (object.alais !== undefined && object.alais !== null) {
            message.alais = object.alais;
        }
        else {
            message.alais = "";
        }
        if (object.verkey !== undefined && object.verkey !== null) {
            message.verkey = object.verkey;
        }
        else {
            message.verkey = "";
        }
        if (object.did !== undefined && object.did !== null) {
            message.did = object.did;
        }
        else {
            message.did = "";
        }
        if (object.role !== undefined && object.role !== null) {
            message.role = object.role;
        }
        else {
            message.role = "";
        }
        return message;
    },
};
var globalThis = (() => {
    if (typeof globalThis !== "undefined")
        return globalThis;
    if (typeof self !== "undefined")
        return self;
    if (typeof window !== "undefined")
        return window;
    if (typeof global !== "undefined")
        return global;
    throw "Unable to locate global object";
})();
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
if (util.Long !== Long) {
    util.Long = Long;
    configure();
}
