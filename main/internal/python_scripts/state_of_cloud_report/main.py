"""
Script that generates a report on the cloud resources outside of Terraform control in the customers
scanned public cloud.
"""
import getopt
import json
import os
import sys
from datetime import datetime

import pandas as pd
from mdutils.mdutils import MdUtils

from helpers.cloud_actor_identificaton import (
    create_markdown_table_cloud_actor_summary,
    process_cloud_actor_actions,
)
from helpers.new_resources_and_cost_estimation import (
    create_markdown_table_cost_summary,
    create_new_resource_tabular_breakdowns_with_cost,
    process_new_resources,
    process_pricing_data,
)
from helpers.managed_resource_drift import (
    create_managed_drift_markdown,
)
from helpers.security_scanning import (
    create_markdown_table_security_scans,
    division_to_security_scan_to_df_dict,
)


def create_markdown_file(job_name: str, markdown_text_output_path):
    """Generate and save a state-of-cloud markdown report"""
    with open("mappings/new-resources-to-documents.json", "r") as json_file:
        new_resources = json.loads(json_file.read())

    with open("mappings/resources-to-cloud-actions.json", "r") as json_file:
        resources_to_cloud_actions = json.loads(json_file.read())

    with open("mappings/division-to-cost-estimates.json", "r") as json_file:
        divisions_to_cost_estimates = json.loads(json_file.read())

    with open("mappings/division-to-security-scan.json", "r") as json_file:
        divisions_to_security_scan = json.loads(json_file.read())

    with open("mappings/drift-resources-differences.json", "r") as json_file:
        managed_drift_list_of_dicts = json.loads(json_file.read())

    if managed_drift_list_of_dicts:
        managed_drift_df = pd.DataFrame(managed_drift_list_of_dicts)
        managed_drift_df["ResourcePath"] = (
            managed_drift_df["ModuleName"]
            + ' (module) "'
            + managed_drift_df["ResourceType"]
            + '" "'
            + managed_drift_df["ResourceName"]
            + '"'
        )
    else:
        managed_drift_df = pd.DataFrame()

    resource_count_dict_of_dfs = {}
    if len(new_resources) > 0:
        resource_count_dict_of_dfs = process_new_resources(new_resources=new_resources)

    if resources_to_cloud_actions:
        actor_action_count_df = process_cloud_actor_actions(
            resources_to_cloud_actions=resources_to_cloud_actions,
        )

    if divisions_to_cost_estimates:
        cost_summary_dict_of_dfs = process_pricing_data(
            divisions_to_cost_estimates=divisions_to_cost_estimates,
            new_resources=new_resources,
        )

    markdown_file = MdUtils(
        file_name=f"{markdown_text_output_path}/report.md",
        title=f"{job_name} - State of Scanned Cloud Resources",
    )

    markdown_file.new_header(level=1, title=f"How to Read this Report", style="atx")
    markdown_file.new_line(
        f"Your job, titled {job_name}, has run. Of the resources "
        "the job scans, at least one resource was identified to have drifted or be outside of Terraform control. "
        "While code has been generated of the Terraform code and corresponding import statements needed to bring these "
        "resources under Terraform control, below you will find a summary of the gaps identified in your "
        "current IaC posture."
    )

    markdown_file.new_header(level=1, title="Identified Security Risks", style="atx")
    if divisions_to_security_scan:
        division_to_security_df_dict = division_to_security_scan_to_df_dict(
            divisions_to_security_scan=divisions_to_security_scan
        )

        markdown_file = create_markdown_table_security_scans(
            markdown_file=markdown_file,
            division_to_security_df_dict=division_to_security_df_dict,
        )
    else:
        markdown_file.new_line("Security scan not run.")

    markdown_file.new_header(
        level=1, title="Calculable Cloud Costs (Monthly)", style="atx"
    )
    if divisions_to_cost_estimates:
        markdown_file = create_markdown_table_cost_summary(
            markdown_file=markdown_file,
            cost_summary_df=cost_summary_dict_of_dfs["cost_summary"],
        )
    else:
        markdown_file.new_line("Cost estimation not run.")

    markdown_file.new_header(
        level=1, title="Resources Outside of Terraform Control", style="atx"
    )

    if len(resource_count_dict_of_dfs) > 0:
        markdown_file = create_new_resource_tabular_breakdowns_with_cost(
            markdown_file=markdown_file,
            resource_count_dict_of_dfs=resource_count_dict_of_dfs,
            cost_by_provider_by_type_df=cost_summary_dict_of_dfs[
                "uncontrolled_cost_by_div_by_type_df"
            ]
            if cost_summary_dict_of_dfs
            else pd.DataFrame(),
        )
    else:
        markdown_file.new_line("No new resources found!")

    markdown_file.new_header(
        level=1, title="Drifted Resources Managed By Terraform", style="atx"
    )
    if not managed_drift_df.empty:
        markdown_file = create_managed_drift_markdown(
            managed_drift_df=managed_drift_df,
            markdown_file=markdown_file,
        )
    else:
        markdown_file.new_line("No controlled resources have drifted!")

    markdown_file.new_header(level=1, title="Root Causes of Drift", style="atx")
    markdown_file.new_header(
        level=2,
        title="Cloud Actors Causing Changes",
        add_table_of_contents="n",
    )
    if resources_to_cloud_actions:
        markdown_file, _ = create_markdown_table_cloud_actor_summary(
            actor_action_count_df=actor_action_count_df,
            markdown_file=markdown_file,
        )
    else:
        markdown_file.new_line("No identified Cloud Actor actions.")

    markdown_file.new_header(level=4, title="Disclaimer", add_table_of_contents="n")
    markdown_file.new_line(
        "*Indicates that a resource's cost is usage based. Since we currently do not infer/have knowledge of usage, "
        "costs may be material although indicated as 0 here."
    )
    markdown_file.new_line()
    markdown_file.new_line(
        "This report presents information on the state of your cloud at a point in time and as best Cloud Concierge"
        " is able to determine. Cloud Concierge does not currently scan every cloud resource for every "
        " cloud provider. For a list of supported resources,"
        " please see our [documentation](https://www.docs.dragondrop.cloud/)."
    )
    markdown_file.new_line()
    markdown_file.new_line(
        f"Created by Cloud Concierge at {datetime.now().strftime('%H:%M UTC on %Y-%m-%d')}"
    )

    markdown_file.new_table_of_contents(table_title="Contents", depth=1)
    markdown_file.create_md_file()
    print("Down creating markdown-styled report.")


if __name__ == "__main__":
    argv = sys.argv[1:]

    try:
        opts, _ = getopt.getopt(argv, "j:i:m:", ["job_name=", "job_unique_id="])

        for opt, arg in opts:
            if opt in ["-i", "--job_unique_id"]:
                job_unique_id = arg
            if opt in ["-j", "--job_name"]:
                job_name = arg

        markdown_text_output_path = f"state_of_cloud/"

        os.makedirs(markdown_text_output_path, exist_ok=True)

        create_markdown_file(
            job_name=job_name,
            markdown_text_output_path=markdown_text_output_path,
        )

    except Exception as e:
        raise Exception(f"Error: {e}")
