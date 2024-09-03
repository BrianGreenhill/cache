# rubocop:disable all

Vagrant.configure('2') do |config|
  config.vm.define 'cache' do |cache|
    cache.vm.box = 'ubuntu/focal64'
    cache.vm.hostname = 'cache'
    cache.vm.network 'private_network', ip: '10.0.0.23'
    cache.vm.provision 'shell', path: 'cache/provision-cache.sh'
    cache.vm.provider 'virtualbox' do |v|
      v.memory = 512
      v.cpus = 1
    end
  end

  config.vm.define 'db' do |db|
    db.vm.box = 'ubuntu/focal64'
    db.vm.hostname = 'db'
    db.vm.network 'private_network', ip: '10.0.0.24'
    db.vm.provision 'shell', path: 'db/provision-db.sh'
    db.vm.provision 'docker' do |d|
      d.pull_images 'mysql:8.0'
      d.run 'mysql', args: '-e MYSQL_ALLOW_EMPTY_PASSWORD=yes -e MYSQL_DATABASE=fallback -p 3306:3306'
    end
    db.vm.provider 'virtualbox' do |v|
      v.memory = 2048
      v.cpus = 1
    end
  end

  config.vm.define 'web' do |web|
    web.vm.box = 'ubuntu/focal64'
    web.vm.hostname = 'web'
    web.vm.network 'private_network', ip: '10.0.0.25'
    web.vm.synced_folder 'web/', '/web'
    web.vm.provision 'shell', path: 'web/provision-web.sh'
    web.vm.provision 'docker' do |d|
      d.build_image '/web', args: '-t web:latest'
      d.run 'web:latest'

    end
    web.vm.provider 'virtualbox' do |v|
      v.memory = 512
      v.cpus = 1
    end
  end

  config.vm.define 'obs' do |obs|
    obs.vm.box = 'ubuntu/focal64'
    obs.vm.hostname = 'obs'
    obs.vm.network 'private_network', ip: '10.0.0.26'
    obs.vm.network "forwarded_port", guest: 3100, host: 3100
    obs.vm.network "forwarded_port", guest: 3000, host: 3000
    obs.vm.synced_folder 'obs/', '/obs'
    obs.vm.provision 'docker' do |d|
      d.pull_images 'clickhouse/clickhouse-server:24.7.4'
      d.pull_images 'qxip/qryn:bun'
      d.pull_images 'grafana/grafana-enterprise:latest'

      d.run 'clickhouse',
        image: 'clickhouse/clickhouse-server',
        args: '--net host -e CLICKHOUSE_PASSWORD=pass -e CLICKHOUSE_USER=user'
      d.run 'qryn',
        image: 'qxip/qryn:bun',
        args: '--net host -e CLICKHOUSE_AUTH=user:pass -e CLICKHOUSE_DATABASE=qryn -e NODE_OPTIONS="--max-old-space-size=4096"'
      d.run 'grafana',
        image: 'grafana/grafana-enterprise',
        args: '--net host -e GF_SECURITY_ADMIN_PASSWORD=pass'
    end
    obs.vm.provider 'virtualbox' do |v|
      v.memory = 4096
      v.cpus = 1
    end
  end
end
