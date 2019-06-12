FROM tomcat:7.0.94-jre8

RUN apt-get update 
RUN   apt-get -y install software-properties-common
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys EEA14886
RUN add-apt-repository "deb http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main"
RUN apt-get update && \
    apt-get -y install default-jdk
RUN apt-get install curl unzip

WORKDIR  /

RUN mkdir ovitmp && \
    cd ovitmp && \
    curl https://iweb.dl.sourceforge.net/project/dcm4che/Oviyam/2.7.2/Oviyam-2.7.2-bin.zip > oviyam.zip && \
        unzip oviyam.zip
        
RUN rm -R /usr/local/tomcat/webapps/ROOT/
RUN cp /ovitmp/Oviyam-2.7.2-bin/Oviyam-2.7.2-bin/oviyam2.war /usr/local/tomcat/webapps/ROOT.war
RUN cp /ovitmp/Oviyam-2.7.2-bin/tomcat/*.jar  /usr/local/tomcat/lib
RUN unzip /usr/local/tomcat/webapps/ROOT.war -d  /usr/local/tomcat/webapps/tmp && \
    rm /usr/local/tomcat/webapps/ROOT.war
COPY web.xml /usr/local/tomcat/webapps/tmp/WEB-INF/web.xml
COPY config.html /usr/local/tomcat/webapps/tmp/config.html
COPY server.js /usr/local/tomcat/webapps/tmp/js/server.js
RUN cd /usr/local/tomcat/webapps/tmp && \
    jar -cvf0 /usr/local/tomcat/webapps/ROOT.war * && \
    rm -r /usr/local/tomcat/webapps/tmp


COPY tomcat-users.xml /usr/local/tomcat/conf/tomcat-users.xml
