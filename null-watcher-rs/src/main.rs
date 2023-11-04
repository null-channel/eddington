use futures::{Stream, StreamExt, TryStreamExt, pin_mut};
use kube::{
    api::{Api, ApiResource, DynamicObject, GroupVersionKind, Resource, ResourceExt},
    runtime::{metadata_watcher, watcher, watcher::Event, WatchStreamExt},
};
use serde::de::DeserializeOwned;
use tracing::*;
use serde::{Deserialize, Serialize};
use std::{env, fmt::Debug, collections::HashMap};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let client = kube::Client::try_default().await?;

    // If set will receive only the metadata for watched resources
    let watch_metadata = env::var("WATCH_METADATA").map(|s| s == "1").unwrap_or(false);

    // Take dynamic resource identifiers:
    let group = env::var("GROUP").unwrap_or_else(|_| "kpack.io".into());
    let version = env::var("VERSION").unwrap_or_else(|_| "v1alpha2".into());
    let kind = env::var("KIND").unwrap_or_else(|_| "Image".into());

    // Turn them into a GVK
    let gvk = GroupVersionKind::gvk(&group, &version, &kind);
    // Use API discovery to identify more information about the type (like its plural)
    let (ar, _caps) = kube::discovery::pinned_kind(&client, &gvk).await?;

    // Use the full resource info to create an Api with the ApiResource as its DynamicType
    let api = Api::<DynamicObject>::all_with(client, &ar);
    let wc = watcher::Config::default();

    // Start a metadata or a full resource watch
    let obs = watcher(api, wc).default_backoff().applied_objects();

    pin_mut!(obs);
    let mut cache = HashMap::new();
    while let Some(p) = obs.try_next().await? {
        let key = format!("{}/{}", p.namespace().unwrap(),p.name_any());
        let val = cache.get(&key);
        let new_value: KubeStatus = serde_json::from_value(p.clone().data)?;
        let new_buildstate = BuildState::new(new_value);

        let Some(val) = val else {
            cache.insert(key, new_buildstate);
            continue;
        };
        match val {
            BuildState::Ready(_) => println!("already ready"),
            _ => {
                match new_buildstate.clone() {
                    BuildState::Ready(image) => println!("new build avalible here: {}", image),
                    BuildState::Building => println!("new build building"),
                    BuildState::Error(e) => println!("error!!! {}", e),
                } 
            },
        }

        cache.insert(format!("{}/{}", p.namespace().unwrap(),p.name_any()), new_buildstate);

    }
    Ok(())
}


#[derive(Debug,Serialize, Deserialize,Clone)]
enum BuildState {
    Building,
    Ready(String),
    Error(String),
}

impl BuildState {
    fn new(val: KubeStatus) -> BuildState {
       let thing = &val.status.conditions[0].reason;
       match thing.clone().as_str() { 
           "UpToDate" => Self::Ready(val.status.latest_image.unwrap()),
           "Building" => Self::Building,
            _ => Self::Error(thing.clone()),
       } 
    }
}


#[derive(Debug,Serialize, Deserialize,Clone)]
#[serde(rename_all = "camelCase")]
struct KubeStatus {
    //conditions: Vec<KubeCondition>,
    status: Status,
}

#[derive(Debug,Serialize, Deserialize,Clone)]
#[serde(rename_all = "camelCase")]
struct Status {
    conditions: Vec<KubeCondition>,
    latest_image: Option<String>,
}

#[derive(Debug,Serialize, Deserialize,Clone)]
struct KubeCondition {
    reason: String
}
/*
Status:                                                                                                                                                                                                         │
│   Build Cache Name:  github.com-nullclouds-examples-go-cache                                                                                                                                                    │
│   Build Counter:     8                                                                                                                                                                                          │
│   Conditions:                                                                                                                                                                                                   │
│     Last Transition Time:         2023-11-04T14:30:11Z                                                                                                                                                          │
│     Reason:                       UpToDate                                                                                                                                                                      │
│     Status:                       True                                                                                                                                                                          │
│     Type:                         Ready                                                                                                                                                                         │
│     Last Transition Time:         2023-11-04T14:30:11Z                                                                                                                                                          │
│     Reason:                       BuilderReady                                                                                                                                                                  │
│     Status:                       True                                                                                                                                                                          │
│     Type:                         BuilderReady                                                                                                                                                                  │
│   Latest Build Image Generation:  1                                                                                                                                                                             │
│   Latest Build Reason:            COMMIT                                                                                                                                                                        │
│   Latest Build Ref:               github.com-nullclouds-examples-go-build-8                                                                                                                                     │
│   Latest Image:                   index.docker.io/nullchannel/github.com-nullclouds-examples-go@sha256:cfa1f11d00ab1f5791a9bec45baf01b9ef4e7aadfd1763b0beaded0409bdf258                                         │
│   Latest Stack:                   io.buildpacks.stacks.jammy                                                                                                                                                    │
│   Observed Generation:            1                                                                                                                                                                             │
│ Events:                           <none>
*/
