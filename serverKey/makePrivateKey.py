# Server Private Key Generator

# This program generates a random ed25519 private key for the server.
# it should be copied into the config file
# the public part should be copied into SERVER_PUB_KEY in common/constants.go

# This software is MIT licensed, Copyright 2015 Factom Foundation.

import sys
import platform
import os
import ed25519djb as ed25519
import re
import sys
import base58pkg as base58
import hashlib

ec_public_key_prefix = "592a"

def doublecheck_key_works(privkey, pubkey):
    """
    This function takes in a newly generated private and public key
    It signs a message and then validates the signature with the key.
    if the verification fails, it raises an exception and the script stops.
    """
    message = "message to sign"
    signature = ed25519.signature(message, privkey, pubkey)

    #print "test sig: " + signature.encode('hex')
    #signaturebad = signature[:1] + 'X' + signature[2:]
    #signaturebad2 = signature[:-2] + 'X' + signature[-1:]
    ed25519.checkvalid(signature, message, pubkey)
    # if this script got this far, then the signature validated.
    #print "signature good"

def public_EC_key_to_human(public_key, advanced):
    """
    this function takes a 32 byte ed25519 public key in binary form and
    returns a string with a human readable public key according to this spec:
    https://github.com/FactomProject/FactomDocs/blob/master/factomDataStructureDetails.md#entry-credit-address
    For example.  a binary public key passed like this :
    0000000000000000000000000000000000000000000000000000000000000000
    would return this string:
    EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa
    
    The advanced parameter set to true displays technical internal data
    """

    hex_seed = public_key.encode('hex')
    public_key_prefix = ec_public_key_prefix + hex_seed

    if advanced == True:
        print "public key with prefix:   " + public_key_prefix

    digest = hashlib.sha256(hashlib.sha256(public_key_prefix.decode("hex")).digest()).digest()
    if advanced == True:
        print "public key hash: " + digest.encode('hex')
    checksummed_public_key = public_key_prefix + digest[:4].encode('hex')
    if advanced == True:
        print "public key with checksum: " + checksummed_public_key
    human_pubkey =  base58.b58encode(checksummed_public_key.decode("hex"))
    if advanced == True:
        print "Human readable public key: " + human_pubkey
    return human_pubkey


def main():




    print "Server Private Key Generator"
    print "Server should have run for several minutes prior to this script."
    print "also, before this script runs, mash on the keyboard for a minute or so through the terminal to get more entropy.\n"

    private_key = os.urandom(32)
    pubkey = ed25519.publickey(private_key)
    doublecheck_key_works(private_key, pubkey)

    print "\n Public key for common/constants.go SERVER_PUB_KEY:"
    print pubkey.encode('hex')
    print "\n Private key for ~/.factom/factom.conf :"
    print private_key.encode('hex') + pubkey.encode('hex')
    
    print "\n\n\n"
    
    private_key = os.urandom(32)
    pubkey = ed25519.publickey(private_key)
    doublecheck_key_works(private_key, pubkey)

    print "\n Private EC key for ~/.factom/factom.conf :"
    print private_key.encode('hex') + pubkey.encode('hex')
    print "\n Public Human readable EC address to fund :"
    print public_EC_key_to_human(private_key, False)
    

if "__main__" == __name__:
        main()

