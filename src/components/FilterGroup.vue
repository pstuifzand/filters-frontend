<template>
  <div class="group">
    <div class="panel panel-default conditions">
      <div class="panel-heading" @click="closed = !closed">
        <h3 class="panel-title">Filter: {{ group.name }}</h3>
      </div>

      <div class="panel-body" v-if="!closed">
        <div class="form-group">
          <label for="group-name">Naam</label>
          <input v-model="group.name" class="form-control" id="group-name" placeholder="untitled" />
        </div>

        <div class="panel panel-default conditions">
          <div class="panel-heading">
            <h3 class="panel-title">Voorwaarden</h3>
          </div>

          <div class="panel-body">
            <div class="combine">
              <select v-model="group.combine" class="form-control">
                <option value="and">Alle regels moeten overeenkomen (AND)</option>
                <option value="or">Tenminste &eacute;&eacute;n regel moet overeenkomen (OR)</option>
              </select>
            </div>
      
            <table class="rule-table">
              <tr class="rules" v-for="rule in group.rules">
                <td class="field">
                  <select v-model="rule.field" class="form-control">
                    <option :value="k" v-for="(v, k) in fields">{{ v.name }}</option>
                  </select>
                </td>
                <td class="function">
                  <select v-model="rule.func" class="form-control">
                    <option :value="k" v-for="(v, k) in functions">{{ v.name }}</option>
                  </select>
                </td>
                <td class="arg">
                  <input v-model="rule.arg" class="form-control" />
                </td>
                <td class="remove">
                  <button class="btn btn-danger" @click="remove_rule(rule)" type="button">&cross;</button>
                </td>
              </tr>
              <tr>
                <td colspan="3"></td>
                <td>
                  <button class="btn btn-primary" @click="add_rule" type="button">+</button>
                </td>
              </tr>
            </table>
          </div>
        </div>

        <div class="panel panel-default">
          <div class="panel-heading">
            <h3 class="panel-title">Acties</h3>
          </div>
          <div class="panel-body">
            <table class="rule-table">
              <tr class="rules" v-for="action in group.actions">
                <td class="field">
                  <select v-model="action.action" class="form-control">
                    <option :value="k" v-for="(v, k) in supported_actions">{{ v.name }}</option>
                  </select>
                </td>
                <td class="arg">
                  <input v-model="action.action_value" class="form-control" v-if="need_arg(action.action)" />
                </td>
                <td class="remove">
                  <button class="btn btn-danger" @click="remove_action(action)" type="button">&cross;</button>
                </td>
              </tr>
              <tr>
                <td colspan="2"></td>
                <td>
                  <button class="btn btn-primary" @click="add_action" type="button">+</button>
                </td>
              </tr>
            </table>
          </div>
        </div>

        <button class="btn btn-primary" type="button" @click="$emit('save')">Opslaan</button>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    props: ['group'],

    methods: {
      add_rule () {
        this.group.rules.push({})
      },
      remove_rule (rule) {
        this.group.rules = this.group.rules.filter((v) => v !== rule)
      },
      add_action () {
        this.group.actions.push({func: 'move_to_folder'})
      },
      remove_action (action) {
        this.group.actions = this.group.actions.filter((v) => v !== action)
      },
      need_arg (action) {
        if (action in this.supported_actions) {
          return this.supported_actions[action].arg_count > 0
        }
        return false
      }
    },
    data () {
      return {
        closed: true,

        fields: {
          to: {
            name: 'To'
          },
          from: {
            name: 'From'
          },
          subject: {
            name: 'Subject'
          },
          cc: {
            name: 'Cc'
          },
          bcc: {
            name: 'Bcc'
          },
          body: {
            name: 'Inhoud'
          }
        },
        functions: {
          begins: {
            name: 'Begint met'
          },
          ends: {
            name: 'Eindigt met'
          },
          contains: {
            name: 'Bevat'
          },
          equal: {
            name: 'Gelijk aan'
          },
          not_equal: {
            name: 'Niet gelijk aan'
          }
        },
        supported_actions: {
          move_to_folder: {
            arg_count: 1,
            name: 'Verplaats naar map'
          },
          remove: {
            arg_count: 0,
            name: 'Verwijder bericht'
          }
        }
      }
    }
  }
</script>

<style lang="scss">
.rule-group {
  width: 100%;
  margin-bottom: 12px;

  border: 1px solid #e0e0e0;
  border-radius: 5px;
  padding: 12px;
}
.rule-table {
  border: none;
  width: 100%;

  td {
    padding: 2px;
  }
  td.field {
    width: 20%
  }
  td.function {
    width: 20%
  }
  td.arg {
    width: 100%
  }
  td.remove {
    width: 10%
  }
}
</style>

