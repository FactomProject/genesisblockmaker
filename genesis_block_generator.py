#taken from https://github.com/ethereum/pyethsaletool
#modified to generate the Factom genesis block data.

import bitcoin as b
import sys
import json
import os
import blockchaininfoparser as bipar
import binascii, re
from datetime import datetime

# Timestamp of sale start
start = 1427814000 # mar 31, 2015 10am central. sale start time
# Initial sale rate
initial_rate = 2000
# Initial sale rate duration
initial_period = 7 * 86400
# step size for declining rate phase
rate_decline = 100
# Length of step
rate_period = 7 * 86400
# Number of declining periods
rate_periods = 4
# Final rate
final_rate = 1500
# Period during which final rate is effective
final_period = 10 * 86400 + 3600  # 1h of slack
# Accept post-sale purchases?
post_rate = 0


exodus = '35gLt5EgB367enjSjyEDahhWWcy6p1MGf6' #factoid


minimum = 0
maximum = 150000000000

caches = {}

try:
    os.mkdir('caches')
except:
    pass


# Cache methods that get networking data. Important since this script takes
# a very long time, and will almost certainly be interrupted multiple times
# while in progress
def cache_method_factory(method, filename):
    def new_method(arg):
        if filename not in caches:
            try:
                caches[filename] = json.load(open(filename, 'r'))
            except:
                caches[filename] = {}
        c = caches[filename]
        if str(arg) not in c:
            c[str(arg)] = method(arg)
            json.dump(c, open(filename, 'w'))
        return c[str(arg)]
    return new_method

# Cached versions of the BCI methods that we need
get_block_header_data = cache_method_factory(bipar.get_block_header_data,
                                             'caches/blockheaders.json')
fetchtx = cache_method_factory(bipar.fetchtx, 'caches/fetchtx.json')
history = cache_method_factory(bipar.history, 'caches/history.json')


# Get a dictionary of the transactions and block heights from the history
def get_txs_and_heights(outs):
    txs = {}
    heights = {}
    for o in outs:
        if o['output'][65:] == '0':
            h = o['output'][:64]
            txs[h] = fetchtx(h)
            heights[h] = o['block_height']
            # print txs[h]
            if len(txs) % 50 == 0:
                sys.stderr.write('Processed transactions: %d\n' % len(txs))
    return {"txs": txs, "heights": heights}


# Produce a json list of purchases
def list_purchases(obj):
    txs, heights = obj['txs'], obj['heights']
    o = []
    for h in txs:
        txhex = str(txs[h])
        # print txhex
        txouts = b.deserialize(txhex)['outs']
        if len(txouts) >= 2:
            addr = b.script_to_address(txouts[0]['script'])
            if addr == exodus:
                v = txouts[0]['value'] + 10000  # add 0.1 millibit for the fees paid transfering from koinify wallet to Factom multisig
                ht = heights[h]
                # We care about the timestamp of the previous
                # confirmed block before a transaction
                t = get_block_header_data(ht - 1)['timestamp']
                o.append({
                    "tx": h,
                    "addr": extract_pubkey(txouts[1]['script']),
                    "value": v,
                    "time": t
                })
                if len(o) % 50 == 0:
                    sys.stderr.write('Gathered outputs: %d\n' % len(o))
    return o


def extract_pubkey(script):
    if re.match('^[0-9a-fA-F]*$', script):
        script = binascii.unhexlify(script)
    if script[:10] != 'j(FACTOM00':
        return 0  # this is an error
    else:
        pubkey = script[10:]
        if len (pubkey) != 32:
            return 0  # this is an error
        else:
            return pubkey.encode("hex")

# Compute factoid value from BTC value
def evaluate_purchases(purchases):
    index = 0

    f1=open('./genesis.csv', 'w')
    f1.write('funding txid,conf time unix,conf time human,ed25519 pubkey,# bitcoins,rate,# factoshis,notes\n')

    for p in purchases:
        if p["time"] < start + initial_period:
            rate = initial_rate
        elif p["time"] < start + initial_period + rate_period * rate_periods:
            pid = (p["time"] - (start + initial_period)) // rate_period + 1
            rate = initial_rate - rate_decline * pid
        elif p["time"] < start + initial_period + rate_period * \
                rate_periods + final_period:
            rate = final_rate
        else:
            rate = post_rate

        f1.write(p['tx'])
        f1.write(',')
        f1.write(str(p['time']))
        f1.write(',')
        f1.write(datetime.fromtimestamp(p['time']).strftime('\"%Y-%b-%d %H:%M:%S\"'))
        f1.write(',')
        f1.write(str(p["addr"]))
        f1.write(',')
        f1.write(str(p["value"]))
        f1.write(',')
        f1.write(str(rate))
        f1.write(',')
        f1.write(str(p["value"] * rate))
        f1.write(',')
        f1.write('\n')

        index += 1

    return "wrote genesis.csv"


def evaluate():
    outs = history(exodus)
    sys.stderr.write('Gathered history: %d\n' % len(outs))
    th = get_txs_and_heights(outs)
    sys.stderr.write('Gathered txs and heights\n')
    p = list_purchases(th)
    sys.stderr.write('Listed purchases\n')
    o = evaluate_purchases(p)
    return o

if __name__ == '__main__':
    print json.dumps(evaluate(), indent=4)
