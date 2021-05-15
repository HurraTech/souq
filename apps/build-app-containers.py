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
            image = service["image"]
            if not image.startswith("%s/" % app_name):
                print("ERROR: Invalid image: '%s' images must start with prefix '%s/'" % (image, app_name))
                exit(1)
            images.add(service["image"])

    for image_tag in images:
        # remove prefix
        image = image_tag.replace("%s/" % app_name, "", 1)
        print("Building %s" % image)
        os.chdir("%s/containers/%s" % (app_name, image))
        os.system("docker buildx build  --platform linux/amd64  -t gcr.io/hurrabuild/%s -o type=docker,dest=- . | gzip > ../%s-amd64.tar.gz " % (image_tag,image))
        os.system("docker buildx build  --platform linux/arm64  -t gcr.io/hurrabuild/%s -o type=docker,dest=- . | gzip > ../%s-arm64.tar.gz" % (image_tag, image))

if __name__ == "__main__":
    main()

