# Key List Generator

# This program generates a list of Factoid pubkeys, addresses, and private keys
# it reads from the file names.txt
# it outputs 2 files, one for public and one for private data.

# This software is MIT licensed, Copyright 2015 Factom Foundation.

import sys
import platform
import os
import ed25519djb as ed25519
import re
import sys
import hashlib
import base58pkg as base58

# the value 0x6478 specifies a private key for the base58 encoding
# it results in a string starting with "Fs" for Factoid Secret
factoid_secret_key_prefix = "6478"
factoid_address_prefix = "5fb1"

pageWidth = 80
kerf = 1

def doublecheck_key_works(privkey, pubkey, advanced):
    """
    This function takes in a newly generated private and public key
    It signs a message and then validates the signature with the key.
    if the verification fails, it raises an exception and the script stops.
    """
    message = "message to sign"
    signature = ed25519.signature(message, privkey, pubkey)
    if advanced == True:
        print "test sig: " + signature.encode('hex')
    #signaturebad = signature[:1] + 'X' + signature[2:]
    #signaturebad2 = signature[:-2] + 'X' + signature[-1:]
    ed25519.checkvalid(signature, message, pubkey)
    if advanced == True:
        # if this script got this far, then the signature validated.
        print "signature good"

def private_key_to_human(private_key, advanced):
    """
    this function takes a 32 byte ed25519 private key in binary form and
    returns a string with a human readable private key according to this spec:
    https://github.com/FactomProject/FactomDocs/blob/master/factomDataStructureDetails.md#factoid-private-keys
    For example.  a binary private key passed like this :
    12fab77add10bcabe1b62b3fe8b167e966e4beee38ccf0062fdd207b5906c841
    would return this string:
    Fs1Ts7PsKMwo4ftCYxQJ3rW4pLiRBXyGEjMrxtHycLu52aDgKGEy

    The advanced parameter set to true displays technical internal data
    """

    hex_seed = private_key.encode('hex')
    private_key_prefixed = factoid_secret_key_prefix + hex_seed

    if advanced == True:
        print "Private key with prefix:   " + private_key_prefixed

    digest = hashlib.sha256(hashlib.sha256(private_key_prefixed.decode("hex")).digest()).digest()
    if advanced == True:
        print "Private key hash: " + digest.encode('hex')
    checksummed_private_key = private_key_prefixed + digest[:4].encode('hex')
    if advanced == True:
        print "Private key with checksum: " + checksummed_private_key
    human_privkey =  base58.b58encode(checksummed_private_key.decode("hex"))
    if advanced == True:
        print "Human readable private key: " + human_privkey
    return human_privkey

def public_key_to_human_address(public_key, advanced):
    hex_pubkey = public_key.encode('hex')
    RCD_hex = "01" + hex_pubkey
    rcdHash = hashlib.sha256(hashlib.sha256(RCD_hex.decode("hex")).digest()).digest().encode("hex")
    checksum = hashlib.sha256(hashlib.sha256((factoid_address_prefix + rcdHash).decode("hex")).digest()).digest()
    humanAddress = base58.b58encode((factoid_address_prefix + rcdHash).decode("hex") + checksum[:4])
    return humanAddress

def private_key_to_human_address(private_key, advanced):
    pubkey = ed25519.publickey(private_key)
    return public_key_to_human_address(pubkey, advanced)


def main():

    names = list()
    notes = list()
    privateKeys = list()
    privateKeysHuman = list()
    publicKeys = list()
    publicKeysHuman = list()

    print "Key List Generator"

    with open("names.txt") as f:
        content = f.readlines()

    for i in range(0, len(content)):
        newKey = os.urandom(32)
        privateKeys.append(newKey)
        c = content[i].strip().rsplit(':', 1)
        names.append(c[0])
        if len(c) > 1:
            notes.append(c[1])
        else:
            notes.append("")
        privateKeysHuman.append(private_key_to_human(privateKeys[i], False))

    #### Write private keys

    f = open("privateKeys.txt", "w")

    for i in range(0, len(content)):
        f.write("--------------------------------------------------------------------------------\n")
        #write index
        index = str(i+1)
        for x in range(pageWidth/2-kerf-len(index)):
            f.write(" ")
        f.write(index)
        for x in range(kerf):
            f.write(" ")
        f.write("|")
        for x in range(kerf):
            f.write(" ")
        f.write(index)
        for x in range(pageWidth/2-kerf-len(index)):
            f.write(" ")
        f.write("\n")

        #write key
        left = privateKeysHuman[i][:26]
        right = privateKeysHuman[i][26:]
        for x in range(pageWidth/2-kerf-len(left)):
            f.write(" ")
        f.write(left)
        for x in range(kerf):
            f.write(" ")
        f.write("|")
        for x in range(kerf):
            f.write(" ")
        f.write(right)
        for x in range(pageWidth/2-kerf-len(right)):
            f.write(" ")
        f.write("\n")

        # write name
        for x in range(pageWidth/2-kerf-len(names[i])):
            f.write(" ")
        f.write(names[i])
        for x in range(kerf):
            f.write(" ")
        f.write("|")
        for x in range(kerf):
            f.write(" ")
        f.write(names[i])
        for x in range(pageWidth/2-kerf-len(names[i])):
            f.write(" ")
        f.write("\n")

        # write notes
        chunksize = pageWidth/2 - kerf
        lines = len(notes[i])/(chunksize)

        for y in range(lines+1):
            line = notes[i][y*chunksize:y*chunksize+chunksize]
            for x in range(pageWidth/2-kerf-len(line)):
                f.write(" ")
            f.write(line)
            for x in range(kerf):
                f.write(" ")
            f.write("|")
            for x in range(kerf):
                f.write(" ")
            f.write(line)
            for x in range(pageWidth/2-kerf-len(line)):
                f.write(" ")
            f.write("\n")
    f.close()


    #### Write public keys

    f = open("publicKeys.txt", "w")

    for i in range(0, len(content)):
        publicKeys.append(ed25519.publickey(privateKeys[i]))
        publicKeysHuman.append(public_key_to_human_address(publicKeys[i],False))
        doublecheck_key_works(privateKeys[i], publicKeys[i], False)

        f.write("--------------------------------------------------------------------------------\n")
        #write index
        index = str(i+1)
        f.write(index + "\n")
        #write pubkey human
        f.write(publicKeysHuman[i]+ "\n")
        #write pubkey
        f.write(publicKeys[i].encode("Hex")+ "\n")
        #write name
        f.write(names[i]+ "\n")
        # write notes
        chunksize = pageWidth
        lines = len(notes[i])/(chunksize)

        for y in range(lines+1):
            line = notes[i][y*chunksize:y*chunksize+chunksize]
            f.write(line)
            f.write("\n")
        print names[i]


if "__main__" == __name__:
        main()

