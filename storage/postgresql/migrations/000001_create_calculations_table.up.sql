create table if not exists calculations
(
    id                            bigserial primary key,
    user_id                       bigserial        not null,
    painful_joints                integer          not null,
    swollen_joints                integer          not null,
    physician_activity_assessment integer          not null,
    patient_activity_assessment   integer          not null,
    crp                           double precision not null,
    sdai_index                    text             not null,
    created_at                    timestamptz      not null default now()
);

create index if not exists idx_user_id on calculations (user_id);


