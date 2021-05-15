### Building Apps

```
docker buildx docker buildx create --name mybuilder --use;
cd apps
make openvpn.tar.gz
make syncthing.tar.gz
```