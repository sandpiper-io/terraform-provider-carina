package carina

import (
	"fmt"
	"time"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/getcarina/libcarina"
)

func resourceCarinaCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceCarinaClusterCreate,
		Read:   resourceCarinaClusterRead,
		Update: resourceCarinaClusterUpdate,
		Delete: resourceCarinaClusterDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nodes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 1,
			},
			"autoscale": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"docker_host": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"docker_cert_path": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCarinaClusterCreate(d *schema.ResourceData, meta interface{}) error {
	log.Print("IN CREATE")
	clusterName := d.Get("name").(string)
	nodes := libcarina.Number(d.Get("nodes").(int))
	autoscale := d.Get("autoscale").(bool)
	client := meta.(*libcarina.ClusterClient)

	cluster, err := client.Create(libcarina.Cluster{
		ClusterName: clusterName,
		Nodes:       nodes,
		AutoScale:   autoscale,
	})
	if err != nil {
		log.Print(err)
	}
	log.Print(cluster.Status)
	d.SetId(clusterName)

	for cluster.Status == "new" || cluster.Status == "building" {
		time.Sleep(10 * time.Second)
		cluster, err = client.Get(clusterName)
		if err != nil {
			break
		}
	}
	d.SetId(clusterName)
	d.Set("nodes", cluster.Nodes)
	d.Set("autoscale", cluster.AutoScale)

	host, path, err := DownloadCredentials(client, clusterName)
	if err != nil {
		return err
	}
	d.Set("docker_host", host)
	d.Set("docker_cert_path", path)
	return nil
}

func resourceCarinaClusterRead(d *schema.ResourceData, meta interface{}) error {
	clusterName := d.Get("name").(string)
	client := meta.(*libcarina.ClusterClient)

	cluster, err := client.Get(clusterName)
	if err != nil {
		d.SetId("")
	} else {
		d.Set("nodes", int(cluster.Nodes))
		d.Set("autoscale", cluster.AutoScale)
	}
	return nil
}

func resourceCarinaClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterName := d.Get("name").(string)
	client := meta.(*libcarina.ClusterClient)

	if d.HasChange("autoscale") {
		autoScale := d.Get("autoscale").(bool)
		_, err := client.SetAutoScale(clusterName, autoScale)
		if err != nil {
			return fmt.Errorf(
				"Error updating auto_scle for database: %s", err)
		}
	}
	if d.HasChange("nodes") {
		nodes := d.Get("nodes").(int)
		_, err := client.Grow(clusterName, nodes)
		if err != nil {
			return fmt.Errorf(
				"Error growing nodes for database: %s", err)
		}
	}
	return resourceCarinaClusterRead(d, meta)
}

func resourceCarinaClusterDelete(d *schema.ResourceData, meta interface{}) error {
	clusterName := d.Get("name").(string)
	client := meta.(*libcarina.ClusterClient)

	_, err := client.Delete(clusterName)

	for  {
		time.Sleep(2 * time.Second)
		_, err = client.Get(clusterName)
		if err == nil {
			break
		}
	}

	return nil
}

func DownloadCredentials(client *libcarina.ClusterClient, clusterName string) (string, string, error) {
	credentials, err := client.GetCredentials(clusterName)
	if err != nil {
		return "", "", err
	}

	path, err := CarinaStoreClusterCredentials(credentials, client.Username, clusterName)
	if err != nil {
		return "", "", err
	}

	return credentials.DockerHost, path, err
}

