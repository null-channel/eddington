package templates

var NginxConf = `worker_processes 1;
daemon off;

error_log stderr;
events { worker_connections 1024; }

http {
  charset utf-8;
  log_format cnb 'NginxLog "$request" $status $body_bytes_sent';
  access_log /dev/stdout cnb;
  default_type application/octet-stream;
  include mime.types;
  sendfile on;

  tcp_nopush on;
  keepalive_timeout 30;
  port_in_redirect off; # Ensure that redirects don't include the internal container PORT - 8080

  server {
    listen {{port}};
    root public;
    index index.html index.htm Default.htm;
  }
}`
