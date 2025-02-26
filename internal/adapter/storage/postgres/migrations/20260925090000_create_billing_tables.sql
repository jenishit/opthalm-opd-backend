-- +goose Up
CREATE SEQUENCE invoice_no_seq START 1;

CREATE TABLE invoices (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_no     VARCHAR(20) NOT NULL UNIQUE DEFAULT ('INV-' || lpad(nextval('invoice_no_seq')::text, 6, '0')),
    patient_id     UUID NOT NULL REFERENCES patients(id),
    visit_id       UUID REFERENCES visits(id),
    status         VARCHAR(20) NOT NULL DEFAULT 'draft'
                       CHECK (status IN ('draft', 'finalized', 'cancelled')),
    subtotal       NUMERIC(10, 2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(10, 2) NOT NULL DEFAULT 0,
    tax_amount     NUMERIC(10, 2) NOT NULL DEFAULT 0,
    total_amount   NUMERIC(10, 2) NOT NULL DEFAULT 0,
    paid_amount    NUMERIC(10, 2) NOT NULL DEFAULT 0,
    due_amount     NUMERIC(10, 2) NOT NULL DEFAULT 0,
    payment_status VARCHAR(20) NOT NULL DEFAULT 'unpaid'
                       CHECK (payment_status IN ('unpaid', 'partial', 'paid', 'cancelled')),
    created_by     UUID NOT NULL REFERENCES users(id),
    updated_by     UUID NOT NULL REFERENCES users(id),
    created_at     TIMESTAMP NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE invoice_items (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id        UUID NOT NULL REFERENCES invoices(id),
    bundle_id         UUID,
    item_type         VARCHAR(20) NOT NULL
                          CHECK (item_type IN ('frame', 'lens', 'coating', 'contact_lens', 'service', 'other')),
