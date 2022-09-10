import dns.resolver
import dns.reversename
import socket 


def main():
    name = getName()
    if name != 'bootstrap':
        index = [int(s) for s in name.split('_') if s.isdigit()] #extract digit number found in name after _: e.g. 'ethereum-docker-master_eth_1'
        print(index[0])
    elif name == 'bootstrap':
        print(0)
    else: 
        print('*** setMiner.py Failed *** \nERROR: Container\'s name should be either bootstrap or contain one number after _')
        #print(-1)
       


#web3.geth.personal.list_accounts()
def getIP():
    try: 
        node_name = socket.gethostname() 
        ip = socket.gethostbyname(node_name) 
    except: 
        print("Unable to get Hostname and IP")   
    return ip


def getName():
    # Reverse DNS lookup will result in getting the container ID/Hostname.  
    # We need to force resolving via DNS and bypass OS system call
    p_ip = getIP()
    long_name = str(dns.resolver.query(dns.reversename.from_address(p_ip),"PTR")[0]) 
    if long_name == '':
        print('*** setMiner.py Failed *** \nERROR: Could not get the name using DNS')
   # print(long_name) # e.g. bootstrap.ethereum-docker-master_mynet.
    name = long_name.split('.', 1) # split the long name to get the container's name as it appears when using docker ps -- e.g. ['bootstrap', 'ethereum-docker-master_mynet.']
   # print(name[0]) # e.g. bootstrap
   # [int(s) for s in name[0].split('_') if s.isdigit()] # used to get the last digit in the name
    return name[0]

if __name__ == "__main__":
   main()
