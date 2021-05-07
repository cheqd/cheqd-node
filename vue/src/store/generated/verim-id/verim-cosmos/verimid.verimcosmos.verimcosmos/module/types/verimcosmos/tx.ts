/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";

export const protobufPackage = "verimid.verimcosmos.verimcosmos";

/** this line is used by starport scaffolding # proto/tx/message */
export interface MsgCreateNym {
  creator: string;
  alias: string;
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
  alias: string;
  verkey: string;
  did: string;
  role: string;
}

export interface MsgUpdateNymResponse {}

export interface MsgDeleteNym {
  creator: string;
  id: number;
}

export interface MsgDeleteNymResponse {}

const baseMsgCreateNym: object = {
  creator: "",
  alias: "",
  verkey: "",
  did: "",
  role: "",
};

export const MsgCreateNym = {
  encode(message: MsgCreateNym, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.alias !== "") {
      writer.uint32(18).string(message.alias);
    }
    if (message.verkey !== "") {
      writer.uint32(26).string(message.verkey);
    }
    if (message.did !== "") {
      writer.uint32(34).string(message.did);
    }
    if (message.role !== "") {
      writer.uint32(42).string(message.role);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgCreateNym {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgCreateNym } as MsgCreateNym;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.alias = reader.string();
          break;
        case 3:
          message.verkey = reader.string();
          break;
        case 4:
          message.did = reader.string();
          break;
        case 5:
          message.role = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateNym {
    const message = { ...baseMsgCreateNym } as MsgCreateNym;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.alias !== undefined && object.alias !== null) {
      message.alias = String(object.alias);
    } else {
      message.alias = "";
    }
    if (object.verkey !== undefined && object.verkey !== null) {
      message.verkey = String(object.verkey);
    } else {
      message.verkey = "";
    }
    if (object.did !== undefined && object.did !== null) {
      message.did = String(object.did);
    } else {
      message.did = "";
    }
    if (object.role !== undefined && object.role !== null) {
      message.role = String(object.role);
    } else {
      message.role = "";
    }
    return message;
  },

  toJSON(message: MsgCreateNym): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.alias !== undefined && (obj.alias = message.alias);
    message.verkey !== undefined && (obj.verkey = message.verkey);
    message.did !== undefined && (obj.did = message.did);
    message.role !== undefined && (obj.role = message.role);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgCreateNym>): MsgCreateNym {
    const message = { ...baseMsgCreateNym } as MsgCreateNym;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.alias !== undefined && object.alias !== null) {
      message.alias = object.alias;
    } else {
      message.alias = "";
    }
    if (object.verkey !== undefined && object.verkey !== null) {
      message.verkey = object.verkey;
    } else {
      message.verkey = "";
    }
    if (object.did !== undefined && object.did !== null) {
      message.did = object.did;
    } else {
      message.did = "";
    }
    if (object.role !== undefined && object.role !== null) {
      message.role = object.role;
    } else {
      message.role = "";
    }
    return message;
  },
};

const baseMsgCreateNymResponse: object = { id: 0 };

export const MsgCreateNymResponse = {
  encode(
    message: MsgCreateNymResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgCreateNymResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgCreateNymResponse } as MsgCreateNymResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateNymResponse {
    const message = { ...baseMsgCreateNymResponse } as MsgCreateNymResponse;
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    return message;
  },

  toJSON(message: MsgCreateNymResponse): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgCreateNymResponse>): MsgCreateNymResponse {
    const message = { ...baseMsgCreateNymResponse } as MsgCreateNymResponse;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    return message;
  },
};

const baseMsgUpdateNym: object = {
  creator: "",
  id: 0,
  alias: "",
  verkey: "",
  did: "",
  role: "",
};

export const MsgUpdateNym = {
  encode(message: MsgUpdateNym, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
    }
    if (message.alias !== "") {
      writer.uint32(26).string(message.alias);
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

  decode(input: Reader | Uint8Array, length?: number): MsgUpdateNym {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgUpdateNym } as MsgUpdateNym;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.alias = reader.string();
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

  fromJSON(object: any): MsgUpdateNym {
    const message = { ...baseMsgUpdateNym } as MsgUpdateNym;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    if (object.alias !== undefined && object.alias !== null) {
      message.alias = String(object.alias);
    } else {
      message.alias = "";
    }
    if (object.verkey !== undefined && object.verkey !== null) {
      message.verkey = String(object.verkey);
    } else {
      message.verkey = "";
    }
    if (object.did !== undefined && object.did !== null) {
      message.did = String(object.did);
    } else {
      message.did = "";
    }
    if (object.role !== undefined && object.role !== null) {
      message.role = String(object.role);
    } else {
      message.role = "";
    }
    return message;
  },

  toJSON(message: MsgUpdateNym): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.id !== undefined && (obj.id = message.id);
    message.alias !== undefined && (obj.alias = message.alias);
    message.verkey !== undefined && (obj.verkey = message.verkey);
    message.did !== undefined && (obj.did = message.did);
    message.role !== undefined && (obj.role = message.role);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgUpdateNym>): MsgUpdateNym {
    const message = { ...baseMsgUpdateNym } as MsgUpdateNym;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    if (object.alias !== undefined && object.alias !== null) {
      message.alias = object.alias;
    } else {
      message.alias = "";
    }
    if (object.verkey !== undefined && object.verkey !== null) {
      message.verkey = object.verkey;
    } else {
      message.verkey = "";
    }
    if (object.did !== undefined && object.did !== null) {
      message.did = object.did;
    } else {
      message.did = "";
    }
    if (object.role !== undefined && object.role !== null) {
      message.role = object.role;
    } else {
      message.role = "";
    }
    return message;
  },
};

const baseMsgUpdateNymResponse: object = {};

export const MsgUpdateNymResponse = {
  encode(_: MsgUpdateNymResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgUpdateNymResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgUpdateNymResponse } as MsgUpdateNymResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateNymResponse {
    const message = { ...baseMsgUpdateNymResponse } as MsgUpdateNymResponse;
    return message;
  },

  toJSON(_: MsgUpdateNymResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<MsgUpdateNymResponse>): MsgUpdateNymResponse {
    const message = { ...baseMsgUpdateNymResponse } as MsgUpdateNymResponse;
    return message;
  },
};

const baseMsgDeleteNym: object = { creator: "", id: 0 };

export const MsgDeleteNym = {
  encode(message: MsgDeleteNym, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgDeleteNym {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgDeleteNym } as MsgDeleteNym;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgDeleteNym {
    const message = { ...baseMsgDeleteNym } as MsgDeleteNym;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    return message;
  },

  toJSON(message: MsgDeleteNym): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgDeleteNym>): MsgDeleteNym {
    const message = { ...baseMsgDeleteNym } as MsgDeleteNym;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    return message;
  },
};

const baseMsgDeleteNymResponse: object = {};

export const MsgDeleteNymResponse = {
  encode(_: MsgDeleteNymResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgDeleteNymResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgDeleteNymResponse } as MsgDeleteNymResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgDeleteNymResponse {
    const message = { ...baseMsgDeleteNymResponse } as MsgDeleteNymResponse;
    return message;
  },

  toJSON(_: MsgDeleteNymResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<MsgDeleteNymResponse>): MsgDeleteNymResponse {
    const message = { ...baseMsgDeleteNymResponse } as MsgDeleteNymResponse;
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  /** this line is used by starport scaffolding # proto/tx/rpc */
  CreateNym(request: MsgCreateNym): Promise<MsgCreateNymResponse>;
  UpdateNym(request: MsgUpdateNym): Promise<MsgUpdateNymResponse>;
  DeleteNym(request: MsgDeleteNym): Promise<MsgDeleteNymResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  CreateNym(request: MsgCreateNym): Promise<MsgCreateNymResponse> {
    const data = MsgCreateNym.encode(request).finish();
    const promise = this.rpc.request(
      "verimid.verimcosmos.verimcosmos.Msg",
      "CreateNym",
      data
    );
    return promise.then((data) =>
      MsgCreateNymResponse.decode(new Reader(data))
    );
  }

  UpdateNym(request: MsgUpdateNym): Promise<MsgUpdateNymResponse> {
    const data = MsgUpdateNym.encode(request).finish();
    const promise = this.rpc.request(
      "verimid.verimcosmos.verimcosmos.Msg",
      "UpdateNym",
      data
    );
    return promise.then((data) =>
      MsgUpdateNymResponse.decode(new Reader(data))
    );
  }

  DeleteNym(request: MsgDeleteNym): Promise<MsgDeleteNymResponse> {
    const data = MsgDeleteNym.encode(request).finish();
    const promise = this.rpc.request(
      "verimid.verimcosmos.verimcosmos.Msg",
      "DeleteNym",
      data
    );
    return promise.then((data) =>
      MsgDeleteNymResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
