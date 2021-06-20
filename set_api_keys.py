#!/usr/bin/env python3

path = "cmd/wwhere/keys/keys.go"

def main():
    key = input("Enter API Key: ")
    print("Key: " + key)

    outfile = open(path, "w")
    
    print("Writing key to " + path)
    outfile.write("package keys\n")
    outfile.write("const (\n")
    outfile.write("APIKey = \"" + key + "\"\n")
    outfile.write(")\n")

    outfile.close()


main()