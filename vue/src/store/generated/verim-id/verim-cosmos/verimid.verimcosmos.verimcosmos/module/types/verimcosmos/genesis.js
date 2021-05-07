/* eslint-disable */
import { Nym } from "../verimcosmos/nym";
import { Writer, Reader } from "protobufjs/minimal";
export const protobufPackage = "verimid.verimcosmos.verimcosmos";
const baseGenesisState = {};
export const GenesisState = {
    encode(message, writer = Writer.create()) {
        for (const v of message.nymList) {
            Nym.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseGenesisState };
        message.nymList = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.nymList.push(Nym.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseGenesisState };
        message.nymList = [];
        if (object.nymList !== undefined && object.nymList !== null) {
            for (const e of object.nymList) {
                message.nymList.push(Nym.fromJSON(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.nymList) {
            obj.nymList = message.nymList.map((e) => (e ? Nym.toJSON(e) : undefined));
        }
        else {
            obj.nymList = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseGenesisState };
        message.nymList = [];
        if (object.nymList !== undefined && object.nymList !== null) {
            for (const e of object.nymList) {
                message.nymList.push(Nym.fromPartial(e));
            }
        }
        return message;
    },
};
