rm -f emc_linux_64.tgz
scp -i "~/Documents/aws/edgematrix2.pem" emc_linux_64.zip ubuntu@node-oregon-51d1616aef7f6caa.elb.us-west-2.amazonaws.com:/home/ubuntu/install
ssh -i "~/Documents/aws/edgematrix2.pem" ubuntu@node-oregon-51d1616aef7f6caa.elb.us-west-2.amazonaws.com "sh /home/ubuntu/install/tgz.sh"
scp -i "~/Documents/aws/edgematrix2.pem" ubuntu@node-oregon-51d1616aef7f6caa.elb.us-west-2.amazonaws.com:/home/ubuntu/install/emc_linux_64.tgz .