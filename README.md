### Building Apps

```
docker buildx create --name mybuilder --use;
cd apps
make openvpn.tar.gz
make syncthing.tar.gz
```