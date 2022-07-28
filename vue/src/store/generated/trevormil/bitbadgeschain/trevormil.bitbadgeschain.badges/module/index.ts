// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgNewBadge } from "./types/badges/tx";
import { MsgNewSubBadge } from "./types/badges/tx";
import { MsgRequestTransferBadge } from "./types/badges/tx";
import { MsgHandlePendingTransfer } from "./types/badges/tx";
import { MsgTransferBadge } from "./types/badges/tx";


const types = [
  ["/trevormil.bitbadgeschain.badges.MsgNewBadge", MsgNewBadge],
  ["/trevormil.bitbadgeschain.badges.MsgNewSubBadge", MsgNewSubBadge],
  ["/trevormil.bitbadgeschain.badges.MsgRequestTransferBadge", MsgRequestTransferBadge],
  ["/trevormil.bitbadgeschain.badges.MsgHandlePendingTransfer", MsgHandlePendingTransfer],
  ["/trevormil.bitbadgeschain.badges.MsgTransferBadge", MsgTransferBadge],
  
];
export const MissingWalletError = new Error("wallet is required");

export const registry = new Registry(<any>types);

const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
}

interface SignAndBroadcastOptions {
  fee: StdFee,
  memo?: string
}

const txClient = async (wallet: OfflineSigner, { addr: addr }: TxClientOptions = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw MissingWalletError;
  let client;
  if (addr) {
    client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
  }else{
    client = await SigningStargateClient.offline( wallet, { registry });
  }
  const { address } = (await wallet.getAccounts())[0];

  return {
    signAndBroadcast: (msgs: EncodeObject[], { fee, memo }: SignAndBroadcastOptions = {fee: defaultFee, memo: ""}) => client.signAndBroadcast(address, msgs, fee,memo),
    msgNewBadge: (data: MsgNewBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgNewBadge", value: MsgNewBadge.fromPartial( data ) }),
    msgNewSubBadge: (data: MsgNewSubBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgNewSubBadge", value: MsgNewSubBadge.fromPartial( data ) }),
    msgRequestTransferBadge: (data: MsgRequestTransferBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRequestTransferBadge", value: MsgRequestTransferBadge.fromPartial( data ) }),
    msgHandlePendingTransfer: (data: MsgHandlePendingTransfer): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgHandlePendingTransfer", value: MsgHandlePendingTransfer.fromPartial( data ) }),
    msgTransferBadge: (data: MsgTransferBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgTransferBadge", value: MsgTransferBadge.fromPartial( data ) }),
    
  };
};

interface QueryClientOptions {
  addr: string
}

const queryClient = async ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr });
};

export {
  txClient,
  queryClient,
};
