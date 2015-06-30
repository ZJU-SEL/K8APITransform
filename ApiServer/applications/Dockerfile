# using a tomcat7 image
#FROM apm-jre7-tomcat7:v1
FROM baseimage
# sign name
MAINTAINER zju

#copy the war file into the webapps
#ADD /applications/cxy/logtest-0.0.1.war_deploy/logtest.war  /usr/local/tomcat/webapps/
ADD ./newimage /root/apache-tomcat-7.0.42/webapps/

# port 22
EXPOSE 22

# Tomcat 8080
EXPOSE 8080

CMD /root/apache-tomcat-7.0.42/bin/start.sh
