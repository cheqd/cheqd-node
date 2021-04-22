/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { Nym } from "../verimcosmos/nym";
import {
  PageRequest,
  PageResponse,
} from "../cosmos/base/query/v1beta1/pagination";

export const protobufPackage = "verimid.verimcosmos.verimcosmos";

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

const baseQueryGetNymRequest: object = { id: 0 };

export const QueryGetNymRequest = {
  encode(
    message: QueryGetNymRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetNymRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetNymRequest } as QueryGetNymRequest;
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

  fromJSON(object: any): QueryGetNymRequest {
    const message = { ...baseQueryGetNymRequest } as QueryGetNymRequest;
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    return message;
  },

  toJSON(message: QueryGetNymRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryGetNymRequest>): QueryGetNymRequest {
    const message = { ...baseQueryGetNymRequest } as QueryGetNymRequest;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    return message;
  },
};

const baseQueryGetNymResponse: object = {};

export const QueryGetNymResponse = {
  encode(
    message: QueryGetNymResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.Nym !== undefined) {
      Nym.encode(message.Nym, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetNymResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetNymResponse } as QueryGetNymResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Nym = Nym.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetNymResponse {
    const message = { ...baseQueryGetNymResponse } as QueryGetNymResponse;
    if (object.Nym !== undefined && object.Nym !== null) {
      message.Nym = Nym.fromJSON(object.Nym);
    } else {
      message.Nym = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetNymResponse): unknown {
    const obj: any = {};
    message.Nym !== undefined &&
      (obj.Nym = message.Nym ? Nym.toJSON(message.Nym) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryGetNymResponse>): QueryGetNymResponse {
    const message = { ...baseQueryGetNymResponse } as QueryGetNymResponse;
    if (object.Nym !== undefined && object.Nym !== null) {
      message.Nym = Nym.fromPartial(object.Nym);
    } else {
      message.Nym = undefined;
    }
    return message;
  },
};

const baseQueryAllNymRequest: object = {};

export const QueryAllNymRequest = {
  encode(
    message: QueryAllNymRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAllNymRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryAllNymRequest } as QueryAllNymRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllNymRequest {
    const message = { ...baseQueryAllNymRequest } as QueryAllNymRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllNymRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryAllNymRequest>): QueryAllNymRequest {
    const message = { ...baseQueryAllNymRequest } as QueryAllNymRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllNymResponse: object = {};

export const QueryAllNymResponse = {
  encode(
    message: QueryAllNymResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.Nym) {
      Nym.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAllNymResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryAllNymResponse } as QueryAllNymResponse;
    message.Nym = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Nym.push(Nym.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllNymResponse {
    const message = { ...baseQueryAllNymResponse } as QueryAllNymResponse;
    message.Nym = [];
    if (object.Nym !== undefined && object.Nym !== null) {
      for (const e of object.Nym) {
        message.Nym.push(Nym.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllNymResponse): unknown {
    const obj: any = {};
    if (message.Nym) {
      obj.Nym = message.Nym.map((e) => (e ? Nym.toJSON(e) : undefined));
    } else {
      obj.Nym = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryAllNymResponse>): QueryAllNymResponse {
    const message = { ...baseQueryAllNymResponse } as QueryAllNymResponse;
    message.Nym = [];
    if (object.Nym !== undefined && object.Nym !== null) {
      for (const e of object.Nym) {
        message.Nym.push(Nym.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** this line is used by starport scaffolding # 2 */
  Nym(request: QueryGetNymRequest): Promise<QueryGetNymResponse>;
  NymAll(request: QueryAllNymRequest): Promise<QueryAllNymResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  Nym(request: QueryGetNymRequest): Promise<QueryGetNymResponse> {
    const data = QueryGetNymRequest.encode(request).finish();
    const promise = this.rpc.request(
      "verimid.verimcosmos.verimcosmos.Query",
      "Nym",
      data
    );
    return promise.then((data) => QueryGetNymResponse.decode(new Reader(data)));
  }

  NymAll(request: QueryAllNymRequest): Promise<QueryAllNymResponse> {
    const data = QueryAllNymRequest.encode(request).finish();
    const promise = this.rpc.request(
      "verimid.verimcosmos.verimcosmos.Query",
      "NymAll",
      data
    );
    return promise.then((data) => QueryAllNymResponse.decode(new Reader(data)));
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
