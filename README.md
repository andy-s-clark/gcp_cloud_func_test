Playing with GCP Cloud Functions

1. Create a `Pub/Sub Topic`
2. Grant `myprojectName-123456@appspot.gserviceaccount.com` the `Pub/Sub Publisher` role to the topic
3. Create a GCP Source Repository (and push the code)
4. Deploy the code to a Cloud Function

        PROJECT_ID=myprojectName-123456
        gcloud functions deploy AppMention \
        --region=us-central1 --memory=128MB --runtime=go113 --entry-point=AppMention \
        --timeout=10 --trigger-http \
        --source=https://source.developers.google.com/projects/${PROJECT_ID}/repos/functest/moveable-aliases/master/paths/appmention \
        --update-env-vars=PROJECT_ID=${PROJECT_ID},TOPIC_NAME=functest

5. Create a `Pub/Sub Subscription` (to allow viewing the messages)
6. Submit a Slack `app_mention` event POST request