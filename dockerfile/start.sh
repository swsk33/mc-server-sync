#!/bin/sh
if [ "$EULA" = "true" ]; then
	echo "By setting the EULA to TRUE, you are indicating your agreement to Minecraft EULA (https://aka.ms/MinecraftEULA)."
	echo "eula=true" >eula.txt
fi
/minecraft/mc-sync-server -d
exec java -Xmx${JVM_MAX} -Xms${JVM_MIN} -XX:+UseZGC -jar /minecraft/fabric-server.jar
