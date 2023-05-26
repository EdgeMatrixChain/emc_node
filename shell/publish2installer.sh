echo "upload emc_windows_x64.zip..."
scp -i "~/Documents/aws/proxy4mid.pem" emc_windows_x64.zip ec2-user@18.119.106.181:/home/ec2-user/installer/emc_windows_x64.zip
echo "move emc_windows_x64.zip to installer..."
ssh -i "~/Documents/aws/proxy4mid.pem" ec2-user@18.119.106.181 "sudo mv installer/emc_windows_x64.zip /usr/share/nginx/html/installer/"
echo "upload emc_mac.zip..."
scp -i "~/Documents/aws/proxy4mid.pem" emc_mac.zip ec2-user@18.119.106.181:/home/ec2-user/installer/emc_mac.zip
echo "move emc_mac.zip to installer..."
ssh -i "~/Documents/aws/proxy4mid.pem" ec2-user@18.119.106.181 "sudo mv installer/emc_mac.zip /usr/share/nginx/html/installer/"
echo "upload emc_mac_arm64.zip..."
scp -i "~/Documents/aws/proxy4mid.pem" emc_mac_arm64.zip ec2-user@18.119.106.181:/home/ec2-user/installer/emc_mac_arm64.zip
echo "move emc_mac_arm64.zip to installer..."
ssh -i "~/Documents/aws/proxy4mid.pem" ec2-user@18.119.106.181 "sudo mv installer/emc_mac_arm64.zip /usr/share/nginx/html/installer/"
echo "upload emc_linux_64.zip..."
scp -i "~/Documents/aws/proxy4mid.pem" emc_linux_64.zip ec2-user@18.119.106.181:/home/ec2-user/installer/emc_linux_64.zip
echo "tgz emc_linux_64..."
ssh -i "~/Documents/aws/proxy4mid.pem" ec2-user@18.119.106.181 "sh /home/ec2-user/installer/tgz.sh"
echo "move emc_linux_64.tgz to installer..."
ssh -i "~/Documents/aws/proxy4mid.pem" ec2-user@18.119.106.181 "sudo mv installer/emc_linux_64.tgz /usr/share/nginx/html/installer/"