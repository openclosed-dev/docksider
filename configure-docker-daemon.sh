#!/bin/bash

script_name=$(basename "$0")
now=$(date '+%Y%m%d_%H%M%S')
listen_port=2375

if [ "$(id -u)" -ne 0 ]; then
  echo 'Error: this script must be run with sudo as follows:'
  echo "  sudo bash $script_name"
  exit 1
fi

ip_list=($(hostname -I))
if [ "${#ip_list[@]}" -eq 0 ]; then
  echo 'Cannot detect IP address'
  exit 1
fi

listen_address=${ip_list[0]}

echo 'Detected IP address:' $listen_address

echo 'Generating a new unit file at: /etc/systemd/system/docker.service'
sudo sed -r 's|-H fd://\s+||' /usr/lib/systemd/system/docker.service > /etc/systemd/system/docker.service

config_file=/etc/docker/daemon.json
echo 'Generating configuration file at:' $config_file

sudo mkdir -p /etc/docker

if [ -f "$config_file" ]; then
  sudo cp $config_file "$config_file.$now"
fi

sudo cat <<EOF > $config_file
{
  "hosts": ["unix:///var/run/docker.sock", "tcp://$listen_address:$listen_port"]
}
EOF

echo 'Restarting the Docker daemon...'
sudo systemctl daemon-reload
sudo systemctl restart docker.service

echo 'Done.'
