package entclient

import "github.com/SURF-Innovatie/MORIS/ent"

type Provider struct{ cli *ent.Client }

func New(cli *ent.Client) *Provider { return &Provider{cli: cli} }

func (p *Provider) Client() *ent.Client { return p.cli }
