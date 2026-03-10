BEGIN;

-- Add JSONB columns
ALTER TABLE processes ADD COLUMN data JSONB NOT NULL DEFAULT '{}';
ALTER TABLE processes ADD COLUMN results JSONB NOT NULL DEFAULT '{}';

-- Migrate existing data into JSONB (if any rows exist)
UPDATE processes SET
  data = jsonb_build_object(
    'processDescription', COALESCE(process_description, ''),
    'proposer', COALESCE(proposer, ''),
    'area', COALESCE(area, ''),
    'responsibleManager', COALESCE(responsible_manager, ''),
    'department', COALESCE(department, ''),
    'systemsInvolved', COALESCE(systems_involved, 1),
    'processType', COALESCE(process_type, ''),
    'periodicity', COALESCE(periodicity, ''),
    'frequentChanges', COALESCE(frequent_changes, false),
    'technology', COALESCE(technology, ''),
    'implementationCost', COALESCE(implementation_cost, 0),
    'trainingCost', COALESCE(training_cost, 0),
    'maintenanceCost', COALESCE(maintenance_cost, 0),
    'hourlyCost', COALESCE(hourly_cost, 0),
    'timePerActivity', COALESCE(time_per_activity, 0),
    'activitiesPerDay', COALESCE(activities_per_day, 0),
    'workingDaysPerYear', COALESCE(working_days_per_year, 220),
    'currentErrorRate', COALESCE(current_error_rate, 0),
    'postErrorRate', COALESCE(post_error_rate, 0),
    'errorCost', COALESCE(error_cost, 0),
    'productivityFactor', COALESCE(productivity_factor, 2),
    'timeReductionFactor', COALESCE(time_reduction_factor, 50),
    'dataQualityScore', COALESCE(data_quality_score, 3),
    'auditScore', COALESCE(audit_score, 3),
    'customerExperienceScore', COALESCE(customer_experience_score, 3),
    'errorReductionScore', COALESCE(error_reduction_score, 3),
    'standardizationScore', COALESCE(standardization_score, 3),
    'scalabilityScore', COALESCE(scalability_score, 3)
  ),
  results = jsonb_build_object(
    'operationalSavings', COALESCE(operational_savings, 0),
    'errorReductionSavings', COALESCE(error_reduction_savings, 0),
    'productivityBenefit', COALESCE(productivity_benefit, 0),
    'annualSavings', COALESCE(annual_savings, 0),
    'roi', COALESCE(roi, 0),
    'breakEvenMonths', break_even_months,
    'hoursSavedMonthly', COALESCE(hours_saved_monthly, 0),
    'hoursSavedAnnually', COALESCE(hours_saved_annually, 0),
    'impactScore', COALESCE(impact_score, 0)
  );

-- Drop old columns
ALTER TABLE processes
  DROP COLUMN process_description,
  DROP COLUMN proposer,
  DROP COLUMN area,
  DROP COLUMN responsible_manager,
  DROP COLUMN department,
  DROP COLUMN systems_involved,
  DROP COLUMN process_type,
  DROP COLUMN periodicity,
  DROP COLUMN frequent_changes,
  DROP COLUMN technology,
  DROP COLUMN implementation_cost,
  DROP COLUMN training_cost,
  DROP COLUMN maintenance_cost,
  DROP COLUMN hourly_cost,
  DROP COLUMN time_per_activity,
  DROP COLUMN activities_per_day,
  DROP COLUMN working_days_per_year,
  DROP COLUMN current_error_rate,
  DROP COLUMN post_error_rate,
  DROP COLUMN error_cost,
  DROP COLUMN productivity_factor,
  DROP COLUMN time_reduction_factor,
  DROP COLUMN data_quality_score,
  DROP COLUMN audit_score,
  DROP COLUMN customer_experience_score,
  DROP COLUMN error_reduction_score,
  DROP COLUMN standardization_score,
  DROP COLUMN scalability_score,
  DROP COLUMN operational_savings,
  DROP COLUMN error_reduction_savings,
  DROP COLUMN productivity_benefit,
  DROP COLUMN annual_savings,
  DROP COLUMN roi,
  DROP COLUMN break_even_months,
  DROP COLUMN hours_saved_monthly,
  DROP COLUMN hours_saved_annually,
  DROP COLUMN impact_score;

-- GIN index for JSONB queries
CREATE INDEX idx_processes_data ON processes USING GIN (data);
CREATE INDEX idx_processes_results ON processes USING GIN (results);

COMMIT;
