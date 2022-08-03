// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgFreezeAddress } from "./types/badges/tx";
import { MsgTransferManager } from "./types/badges/tx";
import { MsgRevokeBadge } from "./types/badges/tx";
import { MsgUpdateUris } from "./types/badges/tx";
import { MsgSelfDestructBadge } from "./types/badges/tx";
import { MsgTransferBadge } from "./types/badges/tx";
import { MsgRequestTransferBadge } from "./types/badges/tx";
import { MsgNewSubBadge } from "./types/badges/tx";
import { MsgNewBadge } from "./types/badges/tx";
import { MsgSetApproval } from "./types/badges/tx";
import { MsgUpdatePermissions } from "./types/badges/tx";
import { MsgHandlePendingTransfer } from "./types/badges/tx";
import { MsgRequestTransferManager } from "./types/badges/tx";


const types = [
  ["/trevormil.bitbadgeschain.badges.MsgFreezeAddress", MsgFreezeAddress],
  ["/trevormil.bitbadgeschain.badges.MsgTransferManager", MsgTransferManager],
  ["/trevormil.bitbadgeschain.badges.MsgRevokeBadge", MsgRevokeBadge],
  ["/trevormil.bitbadgeschain.badges.MsgUpdateUris", MsgUpdateUris],
  ["/trevormil.bitbadgeschain.badges.MsgSelfDestructBadge", MsgSelfDestructBadge],
  ["/trevormil.bitbadgeschain.badges.MsgTransferBadge", MsgTransferBadge],
  ["/trevormil.bitbadgeschain.badges.MsgRequestTransferBadge", MsgRequestTransferBadge],
  ["/trevormil.bitbadgeschain.badges.MsgNewSubBadge", MsgNewSubBadge],
  ["/trevormil.bitbadgeschain.badges.MsgNewBadge", MsgNewBadge],
  ["/trevormil.bitbadgeschain.badges.MsgSetApproval", MsgSetApproval],
  ["/trevormil.bitbadgeschain.badges.MsgUpdatePermissions", MsgUpdatePermissions],
  ["/trevormil.bitbadgeschain.badges.MsgHandlePendingTransfer", MsgHandlePendingTransfer],
  ["/trevormil.bitbadgeschain.badges.MsgRequestTransferManager", MsgRequestTransferManager],
  
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
    msgFreezeAddress: (data: MsgFreezeAddress): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgFreezeAddress", value: MsgFreezeAddress.fromPartial( data ) }),
    msgTransferManager: (data: MsgTransferManager): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgTransferManager", value: MsgTransferManager.fromPartial( data ) }),
    msgRevokeBadge: (data: MsgRevokeBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRevokeBadge", value: MsgRevokeBadge.fromPartial( data ) }),
    msgUpdateUris: (data: MsgUpdateUris): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgUpdateUris", value: MsgUpdateUris.fromPartial( data ) }),
    msgSelfDestructBadge: (data: MsgSelfDestructBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgSelfDestructBadge", value: MsgSelfDestructBadge.fromPartial( data ) }),
    msgTransferBadge: (data: MsgTransferBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgTransferBadge", value: MsgTransferBadge.fromPartial( data ) }),
    msgRequestTransferBadge: (data: MsgRequestTransferBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRequestTransferBadge", value: MsgRequestTransferBadge.fromPartial( data ) }),
    msgNewSubBadge: (data: MsgNewSubBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgNewSubBadge", value: MsgNewSubBadge.fromPartial( data ) }),
    msgNewBadge: (data: MsgNewBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgNewBadge", value: MsgNewBadge.fromPartial( data ) }),
    msgSetApproval: (data: MsgSetApproval): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgSetApproval", value: MsgSetApproval.fromPartial( data ) }),
    msgUpdatePermissions: (data: MsgUpdatePermissions): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgUpdatePermissions", value: MsgUpdatePermissions.fromPartial( data ) }),
    msgHandlePendingTransfer: (data: MsgHandlePendingTransfer): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgHandlePendingTransfer", value: MsgHandlePendingTransfer.fromPartial( data ) }),
    msgRequestTransferManager: (data: MsgRequestTransferManager): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRequestTransferManager", value: MsgRequestTransferManager.fromPartial( data ) }),
    
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
