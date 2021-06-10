// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgDeleteNym } from "./types/verimcosmos/tx";
import { MsgCreateNym } from "./types/verimcosmos/tx";
import { MsgUpdateNym } from "./types/verimcosmos/tx";
const types = [
    ["/verimid.verimcosmos.verimcosmos.MsgDeleteNym", MsgDeleteNym],
    ["/verimid.verimcosmos.verimcosmos.MsgCreateNym", MsgCreateNym],
    ["/verimid.verimcosmos.verimcosmos.MsgUpdateNym", MsgUpdateNym],
];
const registry = new Registry(types);
const defaultFee = {
    amount: [],
    gas: "200000",
};
const txClient = async (wallet, { addr: addr } = { addr: "http://localhost:26657" }) => {
    if (!wallet)
        throw new Error("wallet is required");
    const client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
    const { address } = (await wallet.getAccounts())[0];
    return {
        signAndBroadcast: (msgs, { fee = defaultFee, memo = null }) => memo ? client.signAndBroadcast(address, msgs, fee, memo) : client.signAndBroadcast(address, msgs, fee),
        msgDeleteNym: (data) => ({ typeUrl: "/verimid.verimcosmos.verimcosmos.MsgDeleteNym", value: data }),
        msgCreateNym: (data) => ({ typeUrl: "/verimid.verimcosmos.verimcosmos.MsgCreateNym", value: data }),
        msgUpdateNym: (data) => ({ typeUrl: "/verimid.verimcosmos.verimcosmos.MsgUpdateNym", value: data }),
    };
};
const queryClient = async ({ addr: addr } = { addr: "http://localhost:1317" }) => {
    return new Api({ baseUrl: addr });
};
export { txClient, queryClient, };
