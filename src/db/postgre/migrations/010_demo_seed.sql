-- Demo seed data for Kairon hackathon presentation

-- Company: TechFlow Solutions (demo company for Admin level 4)
-- Company: DataPulse Inc (second company)

-- Customers for TechFlow
INSERT INTO supportflow.customers (id, name, email, plan, company) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'Ivan Petrov', 'ivan@techflow.io', 'premium', 'TechFlow Solutions'),
    ('a0000000-0000-0000-0000-000000000002', 'Maria Sidorova', 'maria@techflow.io', 'basic', 'TechFlow Solutions'),
    ('a0000000-0000-0000-0000-000000000003', 'Alex Kozlov', 'alex@techflow.io', 'premium', 'TechFlow Solutions'),
    ('a0000000-0000-0000-0000-000000000004', 'Elena Novikova', 'elena@techflow.io', 'free', 'TechFlow Solutions'),
    ('a0000000-0000-0000-0000-000000000005', 'Dmitry Volkov', 'dmitry@techflow.io', 'basic', 'TechFlow Solutions')
ON CONFLICT (id) DO UPDATE SET company = EXCLUDED.company;

-- Customers for DataPulse
INSERT INTO supportflow.customers (id, name, email, plan, company) VALUES
    ('a0000000-0000-0000-0000-000000000006', 'Sarah Chen', 'sarah@datapulse.com', 'premium', 'DataPulse Inc'),
    ('a0000000-0000-0000-0000-000000000007', 'James Wilson', 'james@datapulse.com', 'basic', 'DataPulse Inc'),
    ('a0000000-0000-0000-0000-000000000008', 'Anna Lopez', 'anna@datapulse.com', 'premium', 'DataPulse Inc')
ON CONFLICT (id) DO UPDATE SET company = EXCLUDED.company;

-- Agents for TechFlow
INSERT INTO supportflow.agents (id, name, email, role, is_online, company) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'Olena Kovalenko', 'olena@kairon.ai', 'agent', true, 'TechFlow Solutions'),
    ('b0000000-0000-0000-0000-000000000002', 'Maksym Shevchenko', 'maksym@kairon.ai', 'agent', true, 'TechFlow Solutions'),
    ('b0000000-0000-0000-0000-000000000003', 'Andrii Bondarenko', 'andrii@kairon.ai', 'lead', true, 'TechFlow Solutions')
ON CONFLICT (id) DO UPDATE SET company = EXCLUDED.company, name = EXCLUDED.name;

-- Agents for DataPulse
INSERT INTO supportflow.agents (id, name, email, role, is_online, company) VALUES
    ('b0000000-0000-0000-0000-000000000004', 'Lisa Park', 'lisa@kairon.ai', 'agent', true, 'DataPulse Inc'),
    ('b0000000-0000-0000-0000-000000000005', 'Tom Harris', 'tom@kairon.ai', 'lead', false, 'DataPulse Inc')
ON CONFLICT (id) DO UPDATE SET company = EXCLUDED.company, name = EXCLUDED.name;

-- Tickets for TechFlow (various statuses)
INSERT INTO supportflow.tickets (id, customer_id, subject, channel, status, priority, category, agent_id, ai_summary, company, created_at, updated_at, closed_at) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'Double charge on my account', 'web', 'resolved', 'high', 'billing_dispute', 'b0000000-0000-0000-0000-000000000001', 'Customer reported double charge of $49.99. AI verified billing history, confirmed duplicate transaction, and processed refund automatically. Customer confirmed satisfaction.', 'TechFlow Solutions', now() - interval '2 days', now() - interval '1 day', now() - interval '1 day'),
    ('c0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000002', 'Cannot reset my password', 'telegram', 'resolved', 'medium', 'account_access', 'b0000000-0000-0000-0000-000000000001', 'Customer unable to reset password via email. AI initiated password reset and confirmed email delivery. Issue resolved automatically.', 'TechFlow Solutions', now() - interval '3 days', now() - interval '2 days', now() - interval '2 days'),
    ('c0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000003', 'Upgrade to premium plan', 'web', 'resolved', 'low', 'plan_change', 'b0000000-0000-0000-0000-000000000002', 'Customer requested upgrade from basic to premium. AI changed plan and confirmed new features access.', 'TechFlow Solutions', now() - interval '4 days', now() - interval '3 days', now() - interval '3 days'),
    ('c0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000004', 'App crashes on login', 'web', 'in_progress', 'high', 'technical_issue', 'b0000000-0000-0000-0000-000000000002', NULL, 'TechFlow Solutions', now() - interval '6 hours', now() - interval '2 hours', NULL),
    ('c0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000005', 'Cancel my subscription', 'email', 'open', 'medium', 'cancellation', NULL, NULL, 'TechFlow Solutions', now() - interval '1 hour', now() - interval '1 hour', NULL),
    ('c0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 'Billing statement incorrect', 'web', 'resolved', 'medium', 'billing_dispute', 'b0000000-0000-0000-0000-000000000001', 'Customer reported incorrect billing statement. AI looked up billing records and confirmed discrepancy. Escalated to billing team for correction.', 'TechFlow Solutions', now() - interval '5 days', now() - interval '4 days', now() - interval '4 days'),
    ('c0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000003', 'Feature request: dark mode', 'web', 'closed', 'low', 'general_inquiry', 'b0000000-0000-0000-0000-000000000003', 'Customer requested dark mode feature. AI acknowledged and logged as feature request.', 'TechFlow Solutions', now() - interval '7 days', now() - interval '6 days', now() - interval '6 days'),
    ('c0000000-0000-0000-0000-000000000008', 'a0000000-0000-0000-0000-000000000002', 'Refund for unused month', 'telegram', 'waiting', 'medium', 'billing_dispute', 'b0000000-0000-0000-0000-000000000001', NULL, 'TechFlow Solutions', now() - interval '12 hours', now() - interval '4 hours', NULL),
    ('c0000000-0000-0000-0000-000000000012', 'a0000000-0000-0000-0000-000000000001', 'Payment method not working', 'web', 'resolved', 'high', 'billing_dispute', 'b0000000-0000-0000-0000-000000000002', 'Customer payment method declined. AI helped update billing info and retry charge.', 'TechFlow Solutions', now() - interval '6 days', now() - interval '5 days', now() - interval '5 days'),
    ('c0000000-0000-0000-0000-000000000013', 'a0000000-0000-0000-0000-000000000005', 'How to use API?', 'web', 'resolved', 'low', 'general_inquiry', 'b0000000-0000-0000-0000-000000000003', 'Customer asked about API usage. AI provided documentation link and examples from knowledge base.', 'TechFlow Solutions', now() - interval '8 days', now() - interval '7 days', now() - interval '7 days'),
    ('c0000000-0000-0000-0000-000000000014', 'a0000000-0000-0000-0000-000000000003', 'Downgrade plan to basic', 'email', 'resolved', 'low', 'plan_change', 'b0000000-0000-0000-0000-000000000002', 'Customer requested downgrade. AI processed plan change and confirmed prorated billing.', 'TechFlow Solutions', now() - interval '9 days', now() - interval '8 days', now() - interval '8 days'),
    ('c0000000-0000-0000-0000-000000000015', 'a0000000-0000-0000-0000-000000000004', 'Two-factor auth not working', 'telegram', 'resolved', 'high', 'account_access', 'b0000000-0000-0000-0000-000000000001', 'Customer locked out due to 2FA issue. AI reset 2FA and sent recovery codes via email.', 'TechFlow Solutions', now() - interval '10 days', now() - interval '9 days', now() - interval '9 days')
ON CONFLICT (id) DO UPDATE SET company = EXCLUDED.company, ai_summary = EXCLUDED.ai_summary, status = EXCLUDED.status, agent_id = EXCLUDED.agent_id, closed_at = EXCLUDED.closed_at;

-- Tickets for DataPulse
INSERT INTO supportflow.tickets (id, customer_id, subject, channel, status, priority, category, agent_id, ai_summary, company, created_at, updated_at, closed_at) VALUES
    ('c0000000-0000-0000-0000-000000000009', 'a0000000-0000-0000-0000-000000000006', 'API rate limit too low', 'web', 'in_progress', 'high', 'technical_issue', 'b0000000-0000-0000-0000-000000000004', NULL, 'DataPulse Inc', now() - interval '3 hours', now() - interval '1 hour', NULL),
    ('c0000000-0000-0000-0000-000000000010', 'a0000000-0000-0000-0000-000000000007', 'Charged after cancellation', 'email', 'resolved', 'high', 'billing_dispute', 'b0000000-0000-0000-0000-000000000004', 'Customer charged $29.99 after cancellation. AI verified cancellation date, confirmed erroneous charge, processed full refund.', 'DataPulse Inc', now() - interval '2 days', now() - interval '1 day', now() - interval '1 day'),
    ('c0000000-0000-0000-0000-000000000011', 'a0000000-0000-0000-0000-000000000008', 'Need enterprise plan details', 'web', 'resolved', 'low', 'general_inquiry', 'b0000000-0000-0000-0000-000000000005', 'Customer inquired about enterprise plan. AI provided plan details and pricing from knowledge base.', 'DataPulse Inc', now() - interval '1 day', now() - interval '20 hours', now() - interval '20 hours'),
    ('c0000000-0000-0000-0000-000000000016', 'a0000000-0000-0000-0000-000000000006', 'SSO integration help', 'web', 'resolved', 'medium', 'technical_issue', 'b0000000-0000-0000-0000-000000000004', 'Customer needed help with SAML SSO setup. AI provided configuration guide from knowledge base.', 'DataPulse Inc', now() - interval '4 days', now() - interval '3 days', now() - interval '3 days'),
    ('c0000000-0000-0000-0000-000000000017', 'a0000000-0000-0000-0000-000000000008', 'Invoice discrepancy', 'email', 'resolved', 'medium', 'billing_dispute', 'b0000000-0000-0000-0000-000000000005', 'Invoice showed wrong amount. AI verified and issued corrected invoice.', 'DataPulse Inc', now() - interval '5 days', now() - interval '4 days', now() - interval '4 days')
ON CONFLICT (id) DO UPDATE SET company = EXCLUDED.company, ai_summary = EXCLUDED.ai_summary, status = EXCLUDED.status, agent_id = EXCLUDED.agent_id, closed_at = EXCLUDED.closed_at;

-- Messages for ticket 1 (double charge - resolved)
INSERT INTO supportflow.messages (ticket_id, role, content, created_at) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'customer', 'Hi, I was charged twice for my subscription this month. $49.99 appeared twice on my credit card statement. Please help!', now() - interval '2 days'),
    ('c0000000-0000-0000-0000-000000000001', 'ai', 'I understand your concern about the double charge. Let me look into your billing history right away.', now() - interval '2 days' + interval '5 seconds'),
    ('c0000000-0000-0000-0000-000000000001', 'ai', 'I found the duplicate transaction. I have processed a refund of $49.99 for the extra charge. You should see it reflected in your account within 3-5 business days. Is there anything else I can help with?', now() - interval '2 days' + interval '15 seconds'),
    ('c0000000-0000-0000-0000-000000000001', 'customer', 'Thank you so much! That was really fast. No, thats all I needed.', now() - interval '1 day')
ON CONFLICT DO NOTHING;

-- Messages for ticket 2 (password reset)
INSERT INTO supportflow.messages (ticket_id, role, content, created_at) VALUES
    ('c0000000-0000-0000-0000-000000000002', 'customer', 'I cant log in to my account. The password reset email never arrives.', now() - interval '3 days'),
    ('c0000000-0000-0000-0000-000000000002', 'ai', 'I will help you reset your password. Let me initiate a new password reset for your account.', now() - interval '3 days' + interval '3 seconds'),
    ('c0000000-0000-0000-0000-000000000002', 'ai', 'I have sent a password reset link to your email maria@techflow.io. Please check your inbox and spam folder. The link will expire in 24 hours.', now() - interval '3 days' + interval '10 seconds'),
    ('c0000000-0000-0000-0000-000000000002', 'customer', 'Got it! I can log in now. Thanks!', now() - interval '2 days')
ON CONFLICT DO NOTHING;

-- Messages for ticket 4 (app crashes - in progress)
INSERT INTO supportflow.messages (ticket_id, role, content, created_at) VALUES
    ('c0000000-0000-0000-0000-000000000004', 'customer', 'The app keeps crashing every time I try to log in on my iPhone. It worked fine yesterday.', now() - interval '6 hours'),
    ('c0000000-0000-0000-0000-000000000004', 'ai', 'I am sorry to hear about the crash issue. Let me look up your account details to investigate this further.', now() - interval '6 hours' + interval '4 seconds'),
    ('c0000000-0000-0000-0000-000000000004', 'ai', 'I see you are on the free plan using iOS 18.2. This appears to be a known issue with the latest app update. Our engineering team is working on a fix. In the meantime, please try clearing the app cache in Settings > Apps > Kairon > Clear Cache. I am escalating this to our technical team as high priority.', now() - interval '6 hours' + interval '20 seconds'),
    ('c0000000-0000-0000-0000-000000000004', 'customer', 'I cleared the cache but it still crashes. Please fix this ASAP.', now() - interval '4 hours'),
    ('c0000000-0000-0000-0000-000000000004', 'agent', 'Hi Elena, this is Maksym from our technical team. We have identified the root cause and a hotfix is being deployed now. It should be available within the next hour. I apologize for the inconvenience.', now() - interval '2 hours')
ON CONFLICT DO NOTHING;

-- Messages for ticket 5 (cancellation - open, pending approval)
INSERT INTO supportflow.messages (ticket_id, role, content, created_at) VALUES
    ('c0000000-0000-0000-0000-000000000005', 'customer', 'I want to cancel my subscription. I dont use the service anymore.', now() - interval '1 hour'),
    ('c0000000-0000-0000-0000-000000000005', 'ai', 'I understand you would like to cancel your subscription. Before I process this, I want to make sure you are aware that you will lose access to all your data and premium features. The cancellation requires agent approval for your protection.', now() - interval '1 hour' + interval '5 seconds')
ON CONFLICT DO NOTHING;

-- Actions for ticket 1 (refund - executed)
INSERT INTO supportflow.actions (ticket_id, type, params, status, result, confidence, created_at, executed_at) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'lookup_customer', '{}', 'executed', '{"success":true,"message":"Customer: Ivan Petrov, Plan: premium, Email: ivan@techflow.io"}', 0.95, now() - interval '2 days' + interval '6 seconds', now() - interval '2 days' + interval '6 seconds'),
    ('c0000000-0000-0000-0000-000000000001', 'lookup_billing', '{}', 'executed', '{"success":true,"message":"Found 3 transactions. Last: $49.99 on 2026-03-13, $49.99 on 2026-03-13 (DUPLICATE), $49.99 on 2026-02-13"}', 0.95, now() - interval '2 days' + interval '8 seconds', now() - interval '2 days' + interval '8 seconds'),
    ('c0000000-0000-0000-0000-000000000001', 'refund', '{"amount":49.99,"reason":"Duplicate charge on 2026-03-13"}', 'approved', '{"success":true,"message":"Refund of $49.99 processed. Transaction ID: TXN-RF-20260313-001"}', 0.9, now() - interval '2 days' + interval '10 seconds', now() - interval '1 day')
ON CONFLICT DO NOTHING;

-- Actions for ticket 2 (password reset - executed)
INSERT INTO supportflow.actions (ticket_id, type, params, status, result, confidence, created_at, executed_at) VALUES
    ('c0000000-0000-0000-0000-000000000002', 'lookup_customer', '{}', 'executed', '{"success":true,"message":"Customer: Maria Sidorova, Plan: basic, Email: maria@techflow.io"}', 0.95, now() - interval '3 days' + interval '4 seconds', now() - interval '3 days' + interval '4 seconds'),
    ('c0000000-0000-0000-0000-000000000002', 'reset_password', '{}', 'executed', '{"success":true,"message":"Password reset email sent to maria@techflow.io"}', 0.92, now() - interval '3 days' + interval '8 seconds', now() - interval '3 days' + interval '8 seconds')
ON CONFLICT DO NOTHING;

-- Actions for ticket 3 (plan change - executed)
INSERT INTO supportflow.actions (ticket_id, type, params, status, result, confidence, created_at, executed_at) VALUES
    ('c0000000-0000-0000-0000-000000000003', 'change_plan', '{"new_plan":"premium"}', 'executed', '{"success":true,"message":"Plan changed from basic to premium"}', 0.93, now() - interval '4 days' + interval '6 seconds', now() - interval '4 days' + interval '6 seconds')
ON CONFLICT DO NOTHING;

-- Actions for ticket 4 (escalation)
INSERT INTO supportflow.actions (ticket_id, type, params, status, result, confidence, created_at, executed_at) VALUES
    ('c0000000-0000-0000-0000-000000000004', 'lookup_customer', '{}', 'executed', '{"success":true,"message":"Customer: Elena Novikova, Plan: free, Email: elena@techflow.io"}', 0.95, now() - interval '6 hours' + interval '5 seconds', now() - interval '6 hours' + interval '5 seconds'),
    ('c0000000-0000-0000-0000-000000000004', 'escalate', '{"reason":"Recurring app crash on iOS, known issue","priority":"high"}', 'executed', '{"success":true,"message":"Escalated to senior support with high priority"}', 0.88, now() - interval '6 hours' + interval '15 seconds', now() - interval '6 hours' + interval '15 seconds')
ON CONFLICT DO NOTHING;

-- Actions for ticket 5 (cancellation - PENDING approval)
INSERT INTO supportflow.actions (ticket_id, type, params, status, result, confidence, created_at) VALUES
    ('c0000000-0000-0000-0000-000000000005', 'cancel_subscription', '{}', 'pending', '{"success":false,"message":"Action requires agent approval before execution"}', 0.9, now() - interval '1 hour' + interval '6 seconds')
ON CONFLICT DO NOTHING;

-- Actions for DataPulse ticket 10 (refund)
INSERT INTO supportflow.actions (ticket_id, type, params, status, result, confidence, created_at, executed_at) VALUES
    ('c0000000-0000-0000-0000-000000000010', 'lookup_billing', '{}', 'executed', '{"success":true,"message":"Subscription cancelled on 2026-03-01. Charge of $29.99 on 2026-03-10 is erroneous."}', 0.95, now() - interval '2 days' + interval '4 seconds', now() - interval '2 days' + interval '4 seconds'),
    ('c0000000-0000-0000-0000-000000000010', 'refund', '{"amount":29.99,"reason":"Charged after cancellation"}', 'approved', '{"success":true,"message":"Refund of $29.99 processed. Transaction ID: TXN-RF-20260313-002"}', 0.9, now() - interval '2 days' + interval '8 seconds', now() - interval '1 day')
ON CONFLICT DO NOTHING;

-- AI Analyses
INSERT INTO supportflow.ai_analyses (ticket_id, intent, sentiment, urgency, suggested_tools, reasoning, confidence) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'billing_dispute', 'angry', 'high', ARRAY['lookup_billing', 'refund'], 'Customer reports duplicate charge, urgent financial issue', 0.95),
    ('c0000000-0000-0000-0000-000000000002', 'account_access', 'negative', 'medium', ARRAY['reset_password'], 'Customer locked out of account, needs password reset', 0.92),
    ('c0000000-0000-0000-0000-000000000003', 'plan_change', 'positive', 'low', ARRAY['change_plan'], 'Customer wants to upgrade, positive buying intent', 0.94),
    ('c0000000-0000-0000-0000-000000000004', 'technical_issue', 'angry', 'high', ARRAY['escalate'], 'App crash blocking user access, requires engineering attention', 0.91),
    ('c0000000-0000-0000-0000-000000000005', 'cancellation', 'neutral', 'medium', ARRAY['cancel_subscription'], 'Customer wants to cancel, retention opportunity', 0.88),
    ('c0000000-0000-0000-0000-000000000006', 'billing_dispute', 'negative', 'medium', ARRAY['lookup_billing'], 'Billing discrepancy reported', 0.90),
    ('c0000000-0000-0000-0000-000000000007', 'general_inquiry', 'positive', 'low', ARRAY[]::text[], 'Feature request, not a support issue', 0.87),
    ('c0000000-0000-0000-0000-000000000008', 'billing_dispute', 'neutral', 'medium', ARRAY['lookup_billing', 'refund'], 'Refund request for unused service period', 0.89),
    ('c0000000-0000-0000-0000-000000000009', 'technical_issue', 'negative', 'high', ARRAY['escalate'], 'API limitations impacting customer operations', 0.93),
    ('c0000000-0000-0000-0000-000000000010', 'billing_dispute', 'angry', 'high', ARRAY['lookup_billing', 'refund'], 'Charged after explicit cancellation, urgent refund needed', 0.96),
    ('c0000000-0000-0000-0000-000000000011', 'general_inquiry', 'positive', 'low', ARRAY[]::text[], 'Information request about enterprise plan', 0.85)
ON CONFLICT DO NOTHING;

-- Knowledge Base for TechFlow
INSERT INTO supportflow.knowledge_base (company, question, answer) VALUES
    ('TechFlow Solutions', 'What is the refund policy?', 'We offer full refunds within 30 days of purchase. After 30 days, prorated refunds are available for annual plans. Duplicate charges are always refunded immediately.'),
    ('TechFlow Solutions', 'How to upgrade my plan?', 'You can upgrade your plan anytime from Settings > Subscription. The price difference is prorated. Premium plan includes: unlimited API calls, priority support, custom integrations, and advanced analytics.'),
    ('TechFlow Solutions', 'What are the available plans?', 'Free: 100 API calls/month, basic support. Basic ($19.99/mo): 10,000 API calls, email support. Premium ($49.99/mo): unlimited API calls, priority support, custom integrations.')
ON CONFLICT DO NOTHING;

-- Knowledge Base for DataPulse
INSERT INTO supportflow.knowledge_base (company, question, answer) VALUES
    ('DataPulse Inc', 'What are the API rate limits?', 'Free: 100 req/min, Basic: 1000 req/min, Premium: 10000 req/min, Enterprise: unlimited. Rate limits reset every 60 seconds.'),
    ('DataPulse Inc', 'Do you offer enterprise plans?', 'Yes! Enterprise plan includes: unlimited API calls, dedicated support engineer, custom SLA (99.99% uptime), on-premise deployment option, and SSO integration. Contact sales@datapulse.com for pricing.')
ON CONFLICT DO NOTHING;

-- Users: TechFlow company admin + support staff
INSERT INTO supportflow.users (email, name, password, level, role, company) VALUES
    ('admin@techflow.io', 'Viktor Moroz', 'demo123', 4, 'Admin', 'TechFlow Solutions'),
    ('olena@techflow.io', 'Olena Kovalenko', 'demo123', 1, 'Support', 'TechFlow Solutions'),
    ('maksym@techflow.io', 'Maksym Shevchenko', 'demo123', 2, 'Senior Support', 'TechFlow Solutions'),
    ('andrii@techflow.io', 'Andrii Bondarenko', 'demo123', 3, 'Team Lead', 'TechFlow Solutions')
ON CONFLICT (email) DO NOTHING;

-- Users: DataPulse company admin
INSERT INTO supportflow.users (email, name, password, level, role, company) VALUES
    ('admin@datapulse.com', 'Lisa Park', 'demo123', 4, 'Admin', 'DataPulse Inc')
ON CONFLICT (email) DO NOTHING;
