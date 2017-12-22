package main

import (
	"fmt"
	"io"
	"strings"
)

func filterCombinator(combine string) (int, string) {
	if combine == "and" {
		return 1, "&&"
	}
	return 0, "||"
}

func generateFilterFile(w io.Writer, fg *Filters) {
	fmt.Fprint(w, `
use strict;
use warnings;

use Email::Filter;
use URI;

open my $log_fh, '>>', '/home/peter/test/logs/email-filter.log' or die "Can't open email-filter.log";

sub convertPath {
    my $url = shift;
    my $uri = URI->new($url);
    my $path = $uri->path;
    $path =~ s{/}{.}g;
    my $dir = "/home/peter/test/Maildir/$path/";
	$dir =~ s{//$}{/};
    print {$log_fh} "convertPath($url, $dir);\n";
	return $dir;
}

sub log_and_check {
	my ($string, $out) = @_;
	print {$log_fh} $string . "\n";
	return $out;
}

sub accept_email {
	my ($email, @where) = @_;
	print "Accepting email: " . join(" ", @where) . "\n";
	$email->accept(@where);
}

my $email = Email::Filter->new(emergency => "~/emergency_mbox");

print {$log_fh} 'From: ' . $email->from() . "\n";
print {$log_fh} 'Subject: ' . $email->subject() . "\n";
print {$log_fh} 'To: ' . $email->to() . "\n";
print {$log_fh} 'Cc: ' . $email->cc() . "\n";
print {$log_fh} 'Bcc: ' . $email->bcc() . "\n";

`)

	for _, f := range fg.Filters {
		fmt.Fprintf(w, "\n# Filter - %s\n", f.Name)

		startValue, op := filterCombinator(f.Combine)

		fmt.Fprintf(w, "my $result_%d = %d;\n", f.Id, startValue)
		for _, r := range f.Rules {
			rule := generateRule(r)
			out := fmt.Sprintf("$result_%d %s= log_and_check(q[%s], %s);", f.Id, op, rule, rule)
			fmt.Fprintln(w, out)
		}

		fmt.Fprintf(w, "if ($result_%d) {\n", f.Id)
		for _, action := range f.Actions {
			fmt.Fprintf(w, "    # actions %s %s\n", action.Action, action.ActionValue)
			if action.Action == "move_to_folder" {
				fmt.Fprintf(w, "    accept_email($email, convertPath(q{%s}));\n", cleanArg(action.ActionValue))
			}
		}

		fmt.Fprint(w, "}\n")
	}
	fmt.Fprint(w, "accept_email($email, convertPath(\"\"));\n")
}

func generateField(r FilterRule) string {
	switch r.Field {
	case "to":
		return `$email->to()`
	case "to or cc": // TODO: improve with cc
		return `$email->to()`
	case "from":
		return `$email->from()`
	case "cc":
		return `$email->cc()`
	case "bcc":
		return `$email->bcc()`
	case "subject":
		return `$email->subject()`
	case "body":
		return `$email->body()`
	}
	return fmt.Sprintf("ERROR[%s]", r.Field)
}

func cleanArg(arg string) string {
	r := strings.NewReplacer("'", "\\'")
	return r.Replace(arg)
}

func generateRule(r FilterRule) string {
	switch r.Function {
	case "begins":
		return fmt.Sprintf("do { my $x = quotemeta q{%s}; %s =~ m{^$x} }", cleanArg(r.Arg), generateField(r))
	case "ends":
		//return fmt.Sprintf("%s =~ m{\\Q%s\\E$}", generateField(r), cleanArg(r.Arg))
		return fmt.Sprintf("do { my $x = quotemeta q{%s}; %s =~ m{$x$} }", cleanArg(r.Arg), generateField(r))
	case "contains":
		//return fmt.Sprintf("%s =~ m{\\Q%s\\E}", generateField(r), cleanArg(r.Arg))
		return fmt.Sprintf("do { my $x = quotemeta q{%s}; %s =~ m{$x} }", cleanArg(r.Arg), generateField(r))
	case "equal":
		return fmt.Sprintf("%s eq q{%s}", generateField(r), r.Arg)
	case "not_equal":
		return fmt.Sprintf("%s ne q{%s}", generateField(r), r.Arg)
	default:
		return fmt.Sprintf("unknown function %s\n", r.Function)
	}
}
