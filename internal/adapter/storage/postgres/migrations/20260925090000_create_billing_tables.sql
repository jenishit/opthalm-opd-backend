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
    description       VARCHAR(255) NOT NULL,
    inventory_item_id UUID,
    quantity          INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price        NUMERIC(10, 2) NOT NULL DEFAULT 0,
    discount_amount   NUMERIC(10, 2) NOT NULL DEFAULT 0,
    line_total        NUMERIC(10, 2) NOT NULL DEFAULT 0,
    created_at        TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_invoice_items_invoice_id ON invoice_items (invoice_id);

CREATE TABLE payments (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id    UUID NOT NULL REFERENCES invoices(id),
    amount        NUMERIC(10, 2) NOT NULL CHECK (amount > 0),
    method        VARCHAR(20) NOT NULL
                      CHECK (method IN ('cash', 'card', 'online', 'bank_transfer')),
    reference_no  VARCHAR(100),
    paid_at       TIMESTAMP NOT NULL DEFAULT now(),
    created_by    UUID NOT NULL REFERENCES users(id),
    created_at    TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_invoice_id ON payments (invoice_id);

-- +goose Down
DROP TABLE payments;
DROP TABLE invoice_items;
DROP TABLE invoices;
DROP SEQUENCE invoice_no_seq;
