use strict;
use feature "say";
use Data::Dumper;
use 5.20.2;
use JSON::XS 'encode_json';

use feature "signatures";
no warnings 'experimental::signatures';

my $filename = shift @ARGV;

open my $fh, '<', $filename  or die "can't open file $filename";

my @rules;
my $rule = {
  actions => [],
};

while (<$fh>) {
  s/\s+$//;
  chomp;

  my ($key, $val) = m{^(\w+)="([^"]+)"$};

  if (%$rule && $key eq 'name') {
    push @rules, $rule;
    $rule = {
      actions => [],
    };
  }

  parse_value($rule, $key, $val);
}

my %function_map = (
  is => 'equal',
  'ends with' => 'ends',
  'begins with' => 'begins',
);

my $id = 1;
for my $rule (@rules) {
  $rule->{id} = $id;
  my $condition = $rule->{condition};
  my $type;
  $rule->{rules} = [];
  while ($condition =~ m/(OR|AND) \(([^\)]+)\)/g) {
    die "different condition $type != $1" if $type && $type ne $1;
    $type = $1;
    my ($field, $function, $arg) = split/,/, $2;
    $function = $function_map{$function} || $function;
    push @{$rule->{rules}}, {
      field => $field,
      func => $function,
      arg => $arg,
    };
  }
  delete $rule->{condition};
  $rule->{combine} = lc $type;
}
continue {
  $id++;
}

shift @rules;
print encode_json({filters => \@rules});

sub parse_value($rule, $key, $val) {
  if ($key eq 'action') {
    if ($val eq 'Move to folder') {
      $val = 'move_to_folder';
    } elsif ($val eq 'JunkScore') {
      $val = 'junk_score';
    } else {
      die $val;
    }
    $rule->{action} = {
      action => $val,
    };
    return;
  }
  elsif ($key eq 'actionValue') {
    my $action = $rule->{action};
    $action->{action_value} = $val;
    push @{$rule->{actions}}, $action;
    delete $rule->{action};
    return;
  }
  if ($key eq 'enabled') {
    $val = $val eq 'yes';
  }
  if ($key eq 'type') {
    return;
  }

  $rule->{$key} = $val;
}
