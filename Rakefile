# Get version from git tag (e.g., "v1.0.0" or "v1.0.0-3-g1a2b3c4")
def git_version
  version = `git describe --tags --always --match 'v*' 2>/dev/null`.strip
  version.empty? ? "dev" : version.sub(/^v/, "")
end

@hobby_version = git_version

desc "Show rake description"
task :default do
    puts
    puts "Run 'rake -T' for a list of tasks."
    puts
    puts "Use 'rake build' to build the 'hobby' binary."
    puts
end

desc "Build the 'hobby' binary"
task :build do
    sh "go build -ldflags '-s -w -X main.Version=#{@hobby_version}' -o bin/hobby github.com/boxesandglue/hobby/cmd/hobby"
end

desc "Show version information"
task :showversion do
    puts "hobby version #{@hobby_version}"
end

desc "Clean build artifacts"
task :clean do
    FileUtils.rm_rf("bin")
end
