-- +goose Up
-- seats_used/patient_records_used/locations_used are deliberately NOT
-- columns here — they're computed live from users/patients/clinic_locations
-- at read time, so they can never drift from reality. Limits are the static
-- ceiling a plan grants (NULL = unlimited) and are set by the platform
-- operator alongside plan_name/status.
--
-- payment_method_label is reserved for a future payment-processor
-- integration (Stripe/local PSP/bank QR — not decided yet) to sync a
-- display label like "Visa ending in 4242" into. Nothing writes to it yet.
ALTER TABLE subscriptions ADD COLUMN seats_limit INTEGER;
ALTER TABLE subscriptions ADD COLUMN patient_records_limit INTEGER;
ALTER TABLE subscriptions ADD COLUMN locations_limit INTEGER;
ALTER TABLE subscriptions ADD COLUMN payment_method_label VARCHAR(100);

-- +goose Down
ALTER TABLE subscriptions DROP COLUMN payment_method_label;
ALTER TABLE subscriptions DROP COLUMN locations_limit;
ALTER TABLE subscriptions DROP COLUMN patient_records_limit;
ALTER TABLE subscriptions DROP COLUMN seats_limit;
