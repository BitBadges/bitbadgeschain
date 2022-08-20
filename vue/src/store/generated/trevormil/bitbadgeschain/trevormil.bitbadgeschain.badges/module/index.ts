// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgTransferManager } from "./types/badges/tx";
import { MsgNewBadge } from "./types/badges/tx";
import { MsgUpdatePermissions } from "./types/badges/tx";
import { MsgHandlePendingTransfer } from "./types/badges/tx";
import { MsgFreezeAddress } from "./types/badges/tx";
import { MsgRequestTransferBadge } from "./types/badges/tx";
import { MsgUpdateBytes } from "./types/badges/tx";
import { MsgSelfDestructBadge } from "./types/badges/tx";
import { MsgRevokeBadge } from "./types/badges/tx";
import { MsgSetApproval } from "./types/badges/tx";
import { MsgNewSubBadge } from "./types/badges/tx";
import { MsgRegisterAddresses } from "./types/badges/tx";
import { MsgRequestTransferManager } from "./types/badges/tx";
import { MsgTransferBadge } from "./types/badges/tx";
import { MsgPruneBalances } from "./types/badges/tx";
import { MsgUpdateUris } from "./types/badges/tx";


const types = [
  ["/trevormil.bitbadgeschain.badges.MsgTransferManager", MsgTransferManager],
  ["/trevormil.bitbadgeschain.badges.MsgNewBadge", MsgNewBadge],
  ["/trevormil.bitbadgeschain.badges.MsgUpdatePermissions", MsgUpdatePermissions],
  ["/trevormil.bitbadgeschain.badges.MsgHandlePendingTransfer", MsgHandlePendingTransfer],
  ["/trevormil.bitbadgeschain.badges.MsgFreezeAddress", MsgFreezeAddress],
  ["/trevormil.bitbadgeschain.badges.MsgRequestTransferBadge", MsgRequestTransferBadge],
  ["/trevormil.bitbadgeschain.badges.MsgUpdateBytes", MsgUpdateBytes],
  ["/trevormil.bitbadgeschain.badges.MsgSelfDestructBadge", MsgSelfDestructBadge],
  ["/trevormil.bitbadgeschain.badges.MsgRevokeBadge", MsgRevokeBadge],
  ["/trevormil.bitbadgeschain.badges.MsgSetApproval", MsgSetApproval],
  ["/trevormil.bitbadgeschain.badges.MsgNewSubBadge", MsgNewSubBadge],
  ["/trevormil.bitbadgeschain.badges.MsgRegisterAddresses", MsgRegisterAddresses],
  ["/trevormil.bitbadgeschain.badges.MsgRequestTransferManager", MsgRequestTransferManager],
  ["/trevormil.bitbadgeschain.badges.MsgTransferBadge", MsgTransferBadge],
  ["/trevormil.bitbadgeschain.badges.MsgPruneBalances", MsgPruneBalances],
  ["/trevormil.bitbadgeschain.badges.MsgUpdateUris", MsgUpdateUris],
  
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
    msgTransferManager: (data: MsgTransferManager): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgTransferManager", value: MsgTransferManager.fromPartial( data ) }),
    msgNewBadge: (data: MsgNewBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgNewBadge", value: MsgNewBadge.fromPartial( data ) }),
    msgUpdatePermissions: (data: MsgUpdatePermissions): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgUpdatePermissions", value: MsgUpdatePermissions.fromPartial( data ) }),
    msgHandlePendingTransfer: (data: MsgHandlePendingTransfer): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgHandlePendingTransfer", value: MsgHandlePendingTransfer.fromPartial( data ) }),
    msgFreezeAddress: (data: MsgFreezeAddress): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgFreezeAddress", value: MsgFreezeAddress.fromPartial( data ) }),
    msgRequestTransferBadge: (data: MsgRequestTransferBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRequestTransferBadge", value: MsgRequestTransferBadge.fromPartial( data ) }),
    msgUpdateBytes: (data: MsgUpdateBytes): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgUpdateBytes", value: MsgUpdateBytes.fromPartial( data ) }),
    msgSelfDestructBadge: (data: MsgSelfDestructBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgSelfDestructBadge", value: MsgSelfDestructBadge.fromPartial( data ) }),
    msgRevokeBadge: (data: MsgRevokeBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRevokeBadge", value: MsgRevokeBadge.fromPartial( data ) }),
    msgSetApproval: (data: MsgSetApproval): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgSetApproval", value: MsgSetApproval.fromPartial( data ) }),
    msgNewSubBadge: (data: MsgNewSubBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgNewSubBadge", value: MsgNewSubBadge.fromPartial( data ) }),
    msgRegisterAddresses: (data: MsgRegisterAddresses): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRegisterAddresses", value: MsgRegisterAddresses.fromPartial( data ) }),
    msgRequestTransferManager: (data: MsgRequestTransferManager): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgRequestTransferManager", value: MsgRequestTransferManager.fromPartial( data ) }),
    msgTransferBadge: (data: MsgTransferBadge): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgTransferBadge", value: MsgTransferBadge.fromPartial( data ) }),
    msgPruneBalances: (data: MsgPruneBalances): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgPruneBalances", value: MsgPruneBalances.fromPartial( data ) }),
    msgUpdateUris: (data: MsgUpdateUris): EncodeObject => ({ typeUrl: "/trevormil.bitbadgeschain.badges.MsgUpdateUris", value: MsgUpdateUris.fromPartial( data ) }),
    
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
