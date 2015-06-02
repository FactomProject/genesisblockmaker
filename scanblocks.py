# Factoid genesis block creation script
# Uses some code from ethereum/vbuterin pyethsaletool
# https://github.com/ethereum/pyethsaletool

# pip install websocket-client
# run btcd with btcd --addrindex

import os
import websocket
import pprint
import json
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


def get_transactions():

    sale_address = '35gLt5EgB367enjSjyEDahhWWcy6p1MGf6'

    pp = pprint.PrettyPrinter(indent=4)

    # create connection to btcd client
    # requires python 2.7.9
    ws = websocket.create_connection("wss://localhost:8334/ws", sslopt={'ssl_version': websocket.ssl.PROTOCOL_TLSv1_2,
                                                      'ca_certs': os.path.join(os.path.dirname(__file__), "rpc.cert")})
    # login to the btcd client
    postdata = json.dumps({'version': '1.1',
                           'method': 'authenticate',
                           'params': ["testuser", "ComplexPWHere"],
                           'id': 1})
    ws.send(postdata)
    login_message = ws.recv()
    print "btcd login message:", login_message

    # get all the transactions associated with sale address
    # should be run after starting btcd with "btcd --addrindex"

    # this is how many transactions to request from btcd.  it should be larger than the expected number of transactions
    max_num_of_transactions = 3000
    postdata = json.dumps({'version': '1.1',
                           'method': 'searchrawtransactions',
                           'params': [sale_address, 1, 0, max_num_of_transactions],
                           'id': 2})
    ws.send(postdata)

    encoded_txs_from_btcd = ws.recv()
    txs_from_btcd = json.loads(encoded_txs_from_btcd)
    print "found", len(txs_from_btcd['result']), "transactions"

    if len(txs_from_btcd['result']) >= max_num_of_transactions:
        print "warning, may be missing some transactions.  asked for", max_num_of_transactions, "and got exactly", len(txs_from_btcd['result'])

    f1=open('./scaned_from_btcd.csv', 'w')
    f1.write('index, funding txid, conf time unix, conf time human, near cutoff time, ed25519 pubkey, # bitcoins, rate, # factoshis\n')

    payment_index=1


    for res in txs_from_btcd['result']:
        if len(res['vout']) == 2:
            for outs in res['vout']:
                if outs['n'] == 1 :
                    if outs['scriptPubKey']['asm'].startswith('OP_RETURN'):
                        f1.write(str(payment_index))
                        f1.write(',')
                        f1.write(res['txid'])
                        f1.write(',')
                        f1.write(str(res['blocktime']))
                        f1.write(',')
                        f1.write(datetime.fromtimestamp(res['blocktime']).strftime('%Y-%m-%d %H:%M:%S'))
                        f1.write(',')
                        f1.write(str(near_cutoff(res)))
                        f1.write(',')
                        f1.write(get_ed25519_pubkey(res))
                        f1.write(',')
                        f1.write(str(get_btc_rx(res)))
                        f1.write(',')
                        f1.write(str(purchase_rate(res)))
                        f1.write(',')
                        f1.write(str(factoshis_purchased(res)))


                        f1.write('\n')
                        payment_index += 1


                        pass
        pass


# given the timestamp from the confirmation time and the number of satoshis, calculate the number of Factoshis
def factoshis_purchased(p):

    num_satoshis = round(get_btc_rx(p))

    rate = purchase_rate(p)

    num_factoshis = rate * num_satoshis

    return num_factoshis


def purchase_rate(p):

    if p["blocktime"] < start + initial_period:
        rate = initial_rate
    elif p["blocktime"] < start + initial_period + rate_period * rate_periods:
        period_count = rate_periods
        while period_count >= 1:
            if p["blocktime"] < start + initial_period + rate_period * period_count:
                rate = initial_rate - (rate_decline * period_count)
            period_count -= 1
    elif p["blocktime"] < start + initial_period + rate_period * rate_periods + final_period:
        rate = final_rate
    else:
        rate = post_rate

    return rate

def near_cutoff(p):

    cutoff_time = start + initial_period
    current_period = 0
    end_of_sale = start + initial_period + (rate_period * rate_periods) + final_period
    while cutoff_time < end_of_sale and current_period <= rate_periods:

        if p["blocktime"] > cutoff_time and p["blocktime"] < (cutoff_time + 36000):
            return cutoff_time - p["blocktime"]
            pass
        else:
            return 0
            pass

        current_period += 1
        cutoff_time += rate_period
    return 0

def get_ed25519_pubkey(trans):

    for outs in trans['vout']:
            if outs['n'] == 1:
                opreturn = outs['scriptPubKey']['hex']
                # stripped_opreturn holds the 32 byte value which the user specified as a pubkey derived from the 12 words
                if len(opreturn) != (2+8+32)*2:
                    return "bad_address"
                stripped_opreturn = opreturn[20:] # 20 is the first 8 bytes + the opreturn command + length assuming varint is 1 byte

                return stripped_opreturn

def get_btc_rx(trans):
    for outs in trans['vout']:
            if outs['n'] == 0:
                satoshis = outs['value'] * 100000000
                satoshis += 10000 # add 0.1 millibit for the fees paid transfering from koinify wallet to Factom multisig

                return satoshis
    return 0



def main():
    get_transactions()


if __name__ == '__main__':
    main()