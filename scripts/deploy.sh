#! /bin/bash
set -Eeuo pipefail

ENV=$1

artifact pull workflow run-deploy.yaml || true
artifact pull workflow run-route.yaml || true
artifact pull workflow migrations.yaml

if [ -f "run-deploy.yam" ]; then
	# Deploy images
	gcloud builds submit \
		--project="btv-platform-${ENV}-2" \
		--no-source \
		--substitutions=_API_SERVICE=api-${ENV},_JOBS_SERVICE=background-worker-${ENV},_CMS_SERVICE=directus-${ENV} \
		--config=run-deploy.yaml
fi

# Run migrations
gcloud builds submit \
	--project="btv-platform-${ENV}-2" \
	--no-source \
	--substitutions=_INSTANCE_CONNECTION_NAME=btv-platform-${ENV}-2:europe-west4:main-instance \
	--config=migrations.yaml

if [ -f "run-route.yam" ]; then
	# Route Traffic
	gcloud builds submit \
		--project="btv-platform-${ENV}-2" \
		--no-source \
		--substitutions=_API_SERVICE=api-${ENV},_JOBS_SERVICE=background-worker-${ENV},_CMS_SERVICE=directus-${ENV} \
		--config=run-route.yaml
fi
