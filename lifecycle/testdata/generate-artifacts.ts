import { execFileSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { readdirSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';
const sdk = resolve(process.argv[2]);
const { callTool } = await import(`${sdk}/dist/esm/builder/tools/registry.js`);
const session = await import(`${sdk}/dist/esm/builder/session/sessionState.js`);
const { normalizeTxMessages } = await import(`${sdk}/dist/esm/cli/utils/normalizeMsg.js`);
Date.now = () => 1893456000000;
const dir = import.meta.dir;
const cases = await Bun.file(`${dir}/builder-inputs.json`).json();
for (const c of cases) {
 const result = await callTool(c.tool, c.input);
 if (result.isError) throw Error(result.text);
 await Bun.write(`${dir}/artifacts/${c.id}.json`, JSON.stringify(result.result, null, 2) + '\n');
}
const creator = 'bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430';
const forever = [{start:'1',end:'18446744073709551615'}];
for(const [id,cap,standard] of [['nft','1','NFTs'],['fungible','10','Fungible Tokens']]) {
 session.resetSession(id);
 session.getOrCreateSession(id);
 session.setStandards(id,[standard]);
 session.setValidTokenIds(id,[{start:'1',end:'1'}]);
 session.setInvariants(id,{noCustomOwnershipTimes:true,maxSupplyPerId:cap,noForcefulPostMintTransfers:true,disablePoolCreation:true});
 session.setDefaultBalances(id,{balances:[],incomingApprovals:[],outgoingApprovals:[],autoApproveAllIncomingTransfers:true,autoApproveSelfInitiatedOutgoingTransfers:true,autoApproveSelfInitiatedIncomingTransfers:true,userPermissions:{}});
 session.addApproval(id,{fromListId:'Mint',toListId:'All',initiatedByListId:creator,approvalId:'manager-mint',tokenIds:[{start:'1',end:'1'}],transferTimes:forever,ownershipTimes:forever,version:'0',approvalCriteria:{overridesFromOutgoingApprovals:true,approvalAmounts:{overallApprovalAmount:cap,perToAddressApprovalAmount:'0',perFromAddressApprovalAmount:'0',perInitiatedByAddressApprovalAmount:'0',amountTrackerId:'mint-cap',resetTimeIntervals:{startTime:'0',intervalLength:'0'}}}});
 session.addApproval(id,{fromListId:'!Mint',toListId:'All',initiatedByListId:'All',approvalId:'transferable-approval',tokenIds:[{start:'1',end:'1'}],transferTimes:forever,ownershipTimes:forever,version:'0',approvalCriteria:{requireFromEqualsInitiatedBy:true}});
 const value=session.getCollectionValue(id);
 value.collectionMetadata={uri:'https://example.com/golden-metadata.json',customData:''};
 value.tokenMetadata=[{uri:'https://example.com/golden-metadata.json',customData:'',tokenIds:[{start:'1',end:'1'}]}];
 const transaction=normalizeTxMessages(session.getTransaction(id));
 const message=transaction.messages[0];
 delete message.value._meta;
 await Bun.write(`${dir}/artifacts/${id}.json`,JSON.stringify(message,null,2)+'\n');
}

const digest=createHash('sha256');
function hashBuild(path:string){for(const entry of readdirSync(path,{withFileTypes:true}).sort((a,b)=>a.name.localeCompare(b.name))){const file=`${path}/${entry.name}`;if(entry.isDirectory())hashBuild(file);else{digest.update(file.slice(sdk.length)).update('\0').update(readFileSync(file)).update('\0')}}}
hashBuild(`${sdk}/dist/esm`);
const git=(args:string[])=>execFileSync('git',args,{cwd:sdk,encoding:'utf8'}).trim();
const artifacts=Object.fromEntries(readdirSync(`${dir}/artifacts`).sort().map(file=>[file,createHash('sha256').update(readFileSync(`${dir}/artifacts/${file}`)).digest('hex')]));
await Bun.write(`${dir}/provenance.json`,JSON.stringify({sdkCommit:git(['rev-parse','HEAD']),sdkDirty:git(['status','--porcelain','--untracked-files=no'])!=='',sdkVersion:(await Bun.file(`${sdk}/package.json`).json()).version,distSha256:digest.digest('hex'),fixtureTimeMs:'1893456000000',artifacts},null,2)+'\n');
