-- +goose Up
INSERT INTO overtime_statuses (code, name, is_terminal)
VALUES ('SUBMITTED', 'Submitted', FALSE),
       ('REJECTED_BY_DEV_HEAD', 'Rejected by Dev HEAD', TRUE),
       ('APPROVED_BY_DEV_HEAD', 'Approved by Dev HEAD', FALSE),
       ('VALIDATING', 'Validating', FALSE),
       ('PENDING_EMPLOYEE_CONFIRMATION', 'Waiting for employee confirm', FALSE),
       ('COMPLETED', 'Completed', TRUE),
       ('EXCESS_APPROVAL', 'Waiting for additional approval', FALSE),
       ('EXCESS_REJECTED', 'Additional rejected', TRUE);

-- +goose Down