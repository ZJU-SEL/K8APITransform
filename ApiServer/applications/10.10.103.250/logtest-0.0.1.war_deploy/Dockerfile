# VERSION 0.0.1
# using a tomcat7 image
FROM cxyreg:5000/apm-jre7-tomcat7:v4
# sign name
MAINTAINER zju

#copy the war file into the webapps
ADD /applications/10.10.103.250/logtest-0.0.1.war_deploy/logtest.war  /root/apache-tomcat-7.0.42-v2/webapps/

# port 22
EXPOSE 22

# Tomcat 8080
EXPOSE 8080

CMD /root/apache-tomcat-7.0.42-v2/bin/start.sh