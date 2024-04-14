sudo -i -u postgres psql -c "CREATE DATABASE bannerservice"
sudo -i -u postgres psql -d bannerservice -c "CREATE ROLE bannerservice WITH LOGIN PASSWORD 'pass777'"
