# config valid only for current version of Capistrano
lock '3.6.1'

set :application, 'go_web'
set :application_bin, 'web'
set :repo_name, 'go_web'
set :repo_url, ->{ "git@github.com:pandamako/#{fetch :repo_name}.git" }

set :branch, ->{ current_branch }
set :scm, :git

set :keep_releases, 5
# set :format, :pretty
set :log_level, :info
# set :pty, true

set :linked_dirs, %w{log tmp/pids tmp/cache tmp/sockets}

set :file_permissions_paths, ['tmp/pids', 'log', './']
set :file_permissions_users, [fetch(:application)]
set :file_permissions_groups, ['apps']

def exec_go_command command
  execute "sudo -u #{fetch :application} -H -i zsh -l -c \"cd #{deploy_to}/current && #{command}\""
end

# не удалять! используется в config/deploy/production.rb
def current_branch
  ENV["BRANCH"] || `git rev-parse --abbrev-ref HEAD`.chomp
end

namespace :test do
  task :git do
    on roles(:app), in: :sequence, wait: 5 do
      execute "git ls-remote #{fetch :repo_url}"
    end
  end
end

namespace :deploy do
  desc "Restart memcached to cleanup application cache"
  task :reset_cache do
    on roles(:app) do
      execute "sudo /etc/init.d/memcached restart"
    end
  end

  namespace :file do
    task :lock do
      on roles(:app) do
        execute "touch /tmp/deploy_#{fetch :application}_#{fetch :stage}.lock"
      end
    end

    task :unlock do
      on roles(:app) do
        execute "rm /tmp/deploy_#{fetch :application}_#{fetch :stage}.lock"
      end
    end
  end
end

namespace :go do
  desc 'install dependencies'
  task :dependencies do
    on roles(:app), in: :sequence, wait: 5 do
      exec_go_command 'go get ./...'
    end
  end
  desc 'compile binary from source'
  task :compile do
    on roles(:app), in: :sequence, wait: 5 do
      # execute "sudo /etc/init.d/#{fetch :application}_#{fetch :stage} stop" exec_go_command ""
      exec_go_command 'go build web.go'
    end
  end

  desc 'restart service'
  task :restart do
    on roles(:app), in: :sequence, wait: 5 do
      execute "sudo /etc/init.d/#{fetch :application_bin}_#{fetch :stage} restart"
    end
  end

  desc 'upgrade service'
  task :upgrade do
    on roles(:app), in: :sequence, wait: 5 do
      execute "sudo /etc/init.d/#{fetch :application_bin}_#{fetch :stage} upgrade"
    end
  end
end

after 'deploy:starting', 'deploy:file:lock'
after 'deploy:published', 'deploy:file:unlock'

after 'deploy:finishing', 'deploy:cleanup'

after 'deploy:published', 'deploy:set_permissions:chmod'
after 'deploy:published', 'deploy:set_permissions:chown'
after 'deploy:published', 'deploy:set_permissions:chgrp'
after 'deploy:published', 'go:dependencies'
after 'deploy:published', 'go:compile'
after 'deploy:published', 'go:upgrade'

Airbrussh.configure do |config|
  config.truncate = false
end
