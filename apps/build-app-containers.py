import os
import sys
import yaml

def main():
    if len(sys.argv) < 2:
        print("Usage: build-app APP_NAME")
        exit(-1)

    app_name = sys.argv[1]
    images = set()
    with open("%s/containers/containers.yml" % app_name, 'r') as f:
        containers = yaml.load(f, Loader=yaml.FullLoader)
        for svc in containers["services"]:
            service = containers["services"][svc]
            #TODO: Validations
            images.add(service["image"])

    for image in images:
        print("Building %s" % image)
        os.chdir("%s/containers/%s" % (app_name, image))
        os.system("docker build . -t %s" % image)
        os.system("docker save %s | gzip > ../%s.tar.gz" % (image, image))
if __name__ == "__main__":
    main()

