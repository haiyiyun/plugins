# 中国省市地区

## 数据库备份

1. mongoexport

        mongoexport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="province" --out=province.json

        mongoexport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="city" --out=city.json

        mongoexport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="area" --out=area.json

        mongoexport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="street" --out=street.json

        mongoexport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="village" --out=village.json

1. 使用工具导出js脚本

        使用Navicat导出js脚本: cities.js

## 数据库恢复

1. mongoimport

        mongoimport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="province" --file=province.json

        mongoimport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="city" --file=city.json

        mongoimport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="area" --file=area.json

        mongoimport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="street" --file=street.json

        mongoimport  --uri="mongodb://local-mongo-rs-27017:27017,local-mongo-rs-27018:27018,local-mongo-rs-27019:27019/?replicaSet=rs0&w=majority&wtimeoutMS=5000&readConcernLevel=majority&readPreference=primary" --db="cities" --collection="village" --file=village.json

1. 使用工具导入js脚本

        使用Navicat导入js脚本: cities.js
