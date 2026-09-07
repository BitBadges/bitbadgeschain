import argparse, base64, json, pathlib, urllib.parse, urllib.request

parser = argparse.ArgumentParser(description="Capture public v35 migration prefixes at one fixed height; no writes to the node.")
parser.add_argument("--height", type=int, required=True)
parser.add_argument("--rpc", default="https://rpc.bitbadges.io")
parser.add_argument("--out", type=pathlib.Path, required=True)
args = parser.parse_args()
HEIGHT = args.height
OUT = args.out
OUT.mkdir(parents=True, exist_ok=True)

def fields(data):
    offset = 0
    def varint():
        nonlocal offset
        result = 0
        for shift in range(0, 70, 7):
            value = data[offset]; offset += 1
            result |= (value & 127) << shift
            if value < 128: return result
        raise ValueError('Invalid protobuf varint')
    while offset < len(data):
        tag = varint()
        assert tag & 7 == 2, tag
        size = varint()
        value = data[offset:offset+size]; offset += size
        assert len(value) == size
        yield tag >> 3, value

def scan(module, prefix):
    name = module + '-' + prefix
    path = OUT / (name + '.json')
    query = urllib.parse.urlencode({'path':json.dumps('/store/'+module+'/subspace'),'data':'0x'+prefix,'height':HEIGHT,'prove':'false'})
    if not path.exists():
        with urllib.request.urlopen(args.rpc.rstrip('/')+'/abci_query?'+query, timeout=45) as response:
            path.write_bytes(response.read(100*1024*1024))
    result = json.loads(path.read_text())['result']['response']
    assert result['code'] == 0, result
    assert int(result['height']) == HEIGHT
    pairs = []
    for tag, data in fields(base64.b64decode(result['value'] or '')):
        assert tag == 1
        pair = dict(fields(data)); key = pair[1]; value = pair.get(2,b'')
        assert key.startswith(bytes.fromhex(prefix)), key
        pairs.append({'key':base64.b64encode(key).decode(),'value':base64.b64encode(value).decode()})
    return pairs

stores = {}
for module, prefixes in {
 'tokenization':['01','02','04','06','07','0b','0d','0f','10','11','14','15','16'],
 'ibcratelimit':['02','04','06','705f696263726174656c696d6974'],
 'managersplitter':['01'],
 'poolmanager':['04']
}.items():
    stores[module] = []
    for prefix in prefixes:
        pairs = scan(module,prefix)
        print(module,prefix,len(pairs),flush=True)
        stores[module].extend(pairs)
(OUT/'state.json').write_text(json.dumps({'height':HEIGHT,'stores':stores}))
