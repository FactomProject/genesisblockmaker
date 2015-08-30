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


def doublecheck_key_works(privkey, pubkey):
    """
    This function takes in a newly generated private and public key
    It signs a message and then validates the signature with the key.
    if the verification fails, it raises an exception and the script stops.
    """
    message = "message to sign"
    signature = ed25519.signature(message, privkey, pubkey)

    print "test sig: " + signature.encode('hex')
    #signaturebad = signature[:1] + 'X' + signature[2:]
    #signaturebad2 = signature[:-2] + 'X' + signature[-1:]
    ed25519.checkvalid(signature, message, pubkey)
    # if this script got this far, then the signature validated.
    print "signature good"



def main():


    print "Server Private Key Generator"
    print "Server should have run for several minutes prior to this script."
    print "also, before this script runs, mash on the keyboard for a minute or so through the terminal to get more entropy.\n"

    private_key = os.urandom(32)
    pubkey = ed25519.publickey(private_key)
    doublecheck_key_works(private_key, pubkey)

    print "\n \n Public key for common/constants.go SERVER_PUB_KEY:"
    print pubkey.encode('hex')
    print "\n Private key for ~/.factom/factom.conf :"
    print private_key.encode('hex') + pubkey.encode('hex')

if "__main__" == __name__:
        main()

